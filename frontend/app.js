const API = '/api';
const ROUTES = ['/profile', '/news', '/requests', '/management', '/settings'];
const state = {
  access: localStorage.getItem('access_token') || '',
  refresh: localStorage.getItem('refresh_token') || '',
  register: false,
  profile: null,
  notificationController: null,
  notificationsEnabled: localStorage.getItem('notifications_enabled') === 'true',
};

const ERROR_MESSAGES = {
  'missing authorization header': 'Необходимо войти в систему.',
  'invalid or expired token': 'Сессия истекла. Войдите снова.',
  unauthorized: 'Необходимо войти в систему.',
  forbidden: 'У вас недостаточно прав для этого действия.',
  'profile not found': 'Сначала заполните профиль.',
  'email already exists': 'Пользователь с таким e-mail уже зарегистрирован.',
  'invalid credentials': 'Неверный e-mail или пароль.',
  'invalid admin code': 'Неверный код сотрудника УК.',
  'request not found': 'Заявка не найдена.',
  'news not found': 'Новость не найдена.',
  'title is required': 'Укажите заголовок.',
  'content is required': 'Введите текст.',
  'description is required': 'Добавьте описание проблемы.',
  'house_id is required': 'Выберите дом.',
};

function friendlyError(data, status) {
  const raw = String(data.error || data.message || '').toLowerCase();
  for (const [fragment, message] of Object.entries(ERROR_MESSAGES)) {
    if (raw.includes(fragment)) return message;
  }
  if (status === 400) return 'Проверьте заполнение формы.';
  if (status === 404) return 'Запрошенные данные не найдены.';
  if (status === 409) return 'Такая запись уже существует.';
  if (status >= 500) return 'Сервис временно недоступен. Попробуйте позже.';
  return 'Не удалось выполнить действие. Попробуйте ещё раз.';
}

const $ = (selector, root = document) => root.querySelector(selector);
const page = $('#page');

function tokenPayload() {
  try {
    const part = state.access.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(decodeURIComponent(escape(atob(part))));
  } catch (_) { return {}; }
}

function isAdmin() { return tokenPayload().role === 'admin'; }
function escapeHTML(value = '') {
  return String(value).replace(/[&<>"']/g, char => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[char]);
}
function formatDate(value) {
  if (!value) return '—';
  const date = new Date(value);
  return Number.isNaN(date.valueOf()) ? escapeHTML(value) : date.toLocaleString('ru-RU', { dateStyle: 'medium', timeStyle: 'short' });
}
function statusName(value) {
  return ({ open: 'Открыта', in_progress: 'В работе', done: 'Выполнена', cancelled: 'Отменена' })[value] || value || '—';
}
function requestType(value) {
  return ({ plumber: 'Сантехника', electrician: 'Электрика' })[value] || value || '—';
}
function toast(message, error = false) {
  const element = $('#toast');
  element.textContent = message;
  element.className = `toast show${error ? ' error' : ''}`;
  clearTimeout(toast.timer);
  toast.timer = setTimeout(() => { element.className = 'toast'; }, 3500);
}

async function api(path, options = {}, retry = true) {
  const headers = { ...(options.body ? { 'Content-Type': 'application/json' } : {}), ...options.headers };
  if (state.access) headers.Authorization = `Bearer ${state.access}`;
  let response;
  try { response = await fetch(`${API}${path}`, { ...options, headers }); }
  catch (_) { throw new Error('Сервер недоступен. Проверьте, что Docker Compose запущен.'); }

  if (response.status === 401 && retry && state.refresh && path !== '/auth/refresh') {
    const refreshResponse = await fetch(`${API}/auth/refresh`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ refresh_token: state.refresh }),
    });
    if (refreshResponse.ok) {
      const tokens = await refreshResponse.json();
      state.access = tokens.access_token;
      localStorage.setItem('access_token', state.access);
      return api(path, options, false);
    }
    logout(false);
  }
  const contentType = response.headers.get('content-type') || '';
  const data = contentType.includes('json') ? await response.json() : {};
  if (!response.ok) {
    const error = new Error(friendlyError(data, response.status));
    error.status = response.status;
    throw error;
  }
  return data;
}

function setTokens(tokens) {
  state.access = tokens.access_token || '';
  state.refresh = tokens.refresh_token || '';
  localStorage.setItem('access_token', state.access);
  localStorage.setItem('refresh_token', state.refresh);
}

function showAuth() {
  $('#authView').classList.remove('hidden');
  $('#appView').classList.add('hidden');
  if (location.pathname !== '/login') history.replaceState({}, '', '/login');
}

function showApp() {
  $('#authView').classList.add('hidden');
  $('#appView').classList.remove('hidden');
  document.querySelectorAll('[data-admin]').forEach(item => item.classList.toggle('hidden', !isAdmin()));
  const payload = tokenPayload();
  $('#userRole').textContent = isAdmin() ? 'Сотрудник УК' : 'Житель';
  $('#userName').textContent = payload.email || (isAdmin() ? 'Администратор' : 'Житель');
  $('#userInitial').textContent = ($('#userName').textContent[0] || 'Д').toUpperCase();
  if (state.notificationsEnabled) startNotificationStream();
}

function navigate(path, push = true) {
  if (!state.access) return showAuth();
  if (!ROUTES.includes(path)) path = '/profile';
  if (path === '/management' && !isAdmin()) path = '/profile';
  if (push && location.pathname !== path) history.pushState({}, '', path);
  document.querySelectorAll('[data-route]').forEach(link => link.classList.toggle('active', link.getAttribute('href') === path));
  const meta = {
    '/profile': ['Личный кабинет', 'Мой дом'], '/news': ['Будьте в курсе', 'Новости'],
    '/requests': ['Помощь рядом', 'Заявки'], '/management': ['Для сотрудников', 'Управление'],
    '/settings': ['Личный кабинет', 'Настройки'],
  }[path];
  $('#pageEyebrow').textContent = meta[0]; $('#pageTitle').textContent = meta[1];
  ({ '/profile': renderProfile, '/news': renderNews, '/requests': renderRequests, '/management': renderManagement, '/settings': renderSettings })[path]();
}

async function renderProfile() {
  page.innerHTML = '<div class="empty">Загружаем данные дома…</div>';
  try {
    const [profile, companies] = await Promise.all([
      api('/profile').catch(error => {
        if (error.status === 404) return { full_name: '', phone: '', apartment: '', house_id: '', uk_id: '' };
        throw error;
      }),
      api('/management-companies'),
    ]);
    state.profile = profile;
    let houses = { houses: [] };
    if (profile.uk_id) houses = await api(`/houses?uk_id=${encodeURIComponent(profile.uk_id)}`);
    const companyOptions = (companies.companies || []).map(x => `<option value="${escapeHTML(x.id)}" ${x.id === profile.uk_id ? 'selected' : ''}>${escapeHTML(x.name)}</option>`).join('');
    const houseOptions = (houses.houses || []).map(x => `<option value="${escapeHTML(x.id)}" ${x.id === profile.house_id ? 'selected' : ''}>${escapeHTML(x.name)} — ${escapeHTML(x.address)}</option>`).join('');
    page.innerHTML = `
      <section class="hero-card"><div><p class="eyebrow">Добро пожаловать</p><h2>${escapeHTML(profile.full_name || 'Заполните профиль')}</h2><p>${profile.apartment ? `Квартира ${escapeHTML(profile.apartment)}` : 'Добавьте адрес и квартиру, чтобы видеть новости дома.'}</p></div></section>
      <div class="grid"><form id="profileForm" class="card full"><div class="card-head"><h3>Личные данные</h3></div><div class="form-grid">
        <label>ФИО<input name="full_name" required value="${escapeHTML(profile.full_name)}"></label>
        <label>Телефон<input name="phone" required value="${escapeHTML(profile.phone)}" placeholder="+7 900 000-00-00"></label>
        <label>Управляющая компания<select id="company" name="uk_id"><option value="">Выберите УК</option>${companyOptions}</select></label>
        <label>Дом<select id="house" name="house_id" required><option value="">Выберите дом</option>${houseOptions}</select></label>
        <label>Квартира<input name="apartment" required value="${escapeHTML(profile.apartment)}"></label>
      </div><div class="actions"><button class="button primary" style="width:auto">Сохранить изменения</button></div></form></div>`;
    $('#company').addEventListener('change', loadHouseOptions);
    $('#profileForm').addEventListener('submit', saveProfile);
  } catch (error) { page.innerHTML = `<div class="empty">${escapeHTML(error.message)}</div>`; }
}

async function loadHouseOptions(event) {
  const select = $('#house'); select.innerHTML = '<option value="">Загрузка…</option>';
  try {
    const result = await api(`/houses?uk_id=${encodeURIComponent(event.target.value)}`);
    select.innerHTML = '<option value="">Выберите дом</option>' + (result.houses || []).map(x => `<option value="${escapeHTML(x.id)}">${escapeHTML(x.name)} — ${escapeHTML(x.address)}</option>`).join('');
  } catch (error) { select.innerHTML = '<option value="">Не удалось загрузить дома</option>'; toast(error.message, true); }
}

async function saveProfile(event) {
  event.preventDefault(); const form = new FormData(event.target);
  try {
    await api('/profile', { method: 'PUT', body: JSON.stringify({ full_name: form.get('full_name'), phone: form.get('phone'), apartment: form.get('apartment'), house_id: form.get('house_id') }) });
    toast('Профиль сохранён'); renderProfile();
  } catch (error) { toast(error.message, true); }
}

async function renderNews() {
  page.innerHTML = '<div class="empty">Загружаем новости…</div>';
  try {
    const data = await api('/news'); const items = data.news || [];
    page.innerHTML = `<div class="grid"><section class="card full"><div class="card-head"><h3>${isAdmin() ? 'Новости всех управляющих компаний' : 'Новости вашего дома'}</h3><button id="refreshNews" class="button secondary">Обновить</button></div><div class="list">${items.length ? items.map(item => `<article class="list-item"><div><h4>${escapeHTML(item.title)}</h4><p>${escapeHTML(item.content)}</p><div class="meta"><span class="badge">${escapeHTML(item.author_name || 'Пользователь')}</span><span class="badge gray">${formatDate(item.created_at)}</span></div></div></article>`).join('') : '<div class="empty">Пока новостей нет</div>'}</div></section></div>`;
    $('#refreshNews').onclick = renderNews;
  } catch (error) { page.innerHTML = `<div class="empty">${escapeHTML(error.message)}</div>`; }
}

async function renderRequests() {
  page.innerHTML = '<div class="empty">Загружаем заявки…</div>';
  try {
    const data = await api('/requests');
    page.innerHTML = `<div class="grid"><form id="requestForm" class="card"><h3>Новая заявка</h3>
      <label>Категория<select name="type"><option value="plumber">Сантехника</option><option value="electrician">Электрика</option></select></label>
      <label>Тема<input name="title" required maxlength="150" placeholder="Кратко опишите проблему"></label>
      <label>Описание<textarea name="description" required placeholder="Что произошло и где?"></textarea></label>
      <button class="button primary">Отправить заявку</button></form>
      <section class="card"><div class="card-head"><h3>Мои заявки</h3><button id="refreshRequests" class="button secondary">Обновить</button></div><div class="list">${requestList(data.requests || [])}</div></section>
      <section id="requestDetails" class="card full hidden"></section></div>`;
    $('#requestForm').onsubmit = createRequest; $('#refreshRequests').onclick = renderRequests;
    document.querySelectorAll('[data-request]').forEach(item => item.onclick = () => openRequest(item.dataset.request));
  } catch (error) { page.innerHTML = `<div class="empty">${escapeHTML(error.message)}</div>`; }
}

function requestList(items) {
  if (!items.length) return '<div class="empty">У вас пока нет заявок</div>';
  return items.map(item => `<article class="list-item" data-request="${escapeHTML(item.id)}"><div><h4>${escapeHTML(item.title)}</h4><p>${escapeHTML(item.description)}</p><div class="meta"><span class="badge">${escapeHTML(item.author_name || 'Пользователь')}</span><span class="badge">${escapeHTML(requestType(item.type))}</span><span class="badge ${item.status === 'open' ? 'orange' : 'gray'}">${escapeHTML(statusName(item.status))}</span></div></div><small>${formatDate(item.created_at)}</small></article>`).join('');
}

async function createRequest(event) {
  event.preventDefault(); const form = new FormData(event.target);
  try { await api('/requests', { method: 'POST', body: JSON.stringify(Object.fromEntries(form)) }); toast('Заявка отправлена'); renderRequests(); }
  catch (error) { toast(error.message, true); }
}

async function openRequest(id) {
  try {
    const [item, comments] = await Promise.all([api(`/requests/${id}`), api(`/requests/${id}/comments`)]);
    const details = $('#requestDetails'); details.classList.remove('hidden');
    details.innerHTML = `<div class="card-head"><h3>${escapeHTML(item.title)}</h3><span class="badge">${escapeHTML(statusName(item.status))}</span></div><p><b>${escapeHTML(item.author_name || 'Пользователь')}</b></p><p>${escapeHTML(item.description)}</p><h3>Комментарии</h3><div class="list">${(comments.comments || []).map(x => `<div class="list-item"><div><b>${escapeHTML(x.author_name || 'Пользователь')}</b><p>${escapeHTML(x.content)}</p><small>${formatDate(x.created_at)}</small></div></div>`).join('') || '<div class="empty">Комментариев пока нет</div>'}</div><form id="commentForm" class="actions"><input name="content" required placeholder="Добавить комментарий"><button class="button secondary">Отправить</button></form>`;
    $('#commentForm').onsubmit = async event => { event.preventDefault(); const content = new FormData(event.target).get('content'); try { await api(`/requests/${id}/comment`, { method: 'POST', body: JSON.stringify({ content }) }); toast('Комментарий добавлен'); openRequest(id); } catch (error) { toast(error.message, true); } };
    details.scrollIntoView({ behavior: 'smooth', block: 'start' });
  } catch (error) { toast(error.message, true); }
}

async function renderManagement() {
  if (!isAdmin()) return navigate('/profile');
  page.innerHTML = '<div class="empty">Загружаем панель управления…</div>';
  try {
    const [companies, houses, requests] = await Promise.all([api('/management-companies'), api('/houses'), api('/requests')]);
    const options = (companies.companies || []).map(x => `<option value="${escapeHTML(x.id)}">${escapeHTML(x.name)}</option>`).join('');
    const houseOptions = (houses.houses || []).map(x => `<option value="${escapeHTML(x.id)}">${escapeHTML(x.name)} — ${escapeHTML(x.address)}</option>`).join('');
    page.innerHTML = `<div class="grid">
      <form id="companyForm" class="card third"><h3>Новая УК</h3><label>Название<input name="name" required></label><button class="button secondary">Создать</button></form>
      <form id="houseForm" class="card third"><h3>Новый дом</h3><label>Название<input name="name" required></label><label>Адрес<input name="address" required></label><label>УК<select name="uk_id" required><option value="">Выберите</option>${options}</select></label><button class="button secondary">Добавить</button></form>
      <form id="newsForm" class="card third"><h3>Опубликовать новость</h3><label>Дом любой УК<select name="house_id" required><option value="">Выберите дом</option>${houseOptions}</select></label><label>Заголовок<input name="title" required></label><label>Текст<textarea name="content" required></textarea></label><button class="button secondary">Опубликовать</button></form>
      <section class="card full"><h3>Заявки жителей</h3><div class="list">${(requests.requests || []).map(item => `<div class="list-item"><div><h4>${escapeHTML(item.title)}</h4><p><b>${escapeHTML(item.author_name || 'Пользователь')}</b></p><p>${escapeHTML(item.description)}</p></div><select data-status="${escapeHTML(item.id)}"><option value="open" ${item.status === 'open' ? 'selected' : ''}>Открыта</option><option value="in_progress" ${item.status === 'in_progress' ? 'selected' : ''}>В работе</option><option value="done" ${item.status === 'done' ? 'selected' : ''}>Выполнена</option><option value="cancelled" ${item.status === 'cancelled' ? 'selected' : ''}>Отменена</option></select></div>`).join('') || '<div class="empty">Заявок нет</div>'}</div></section></div>`;
    bindCreateForm('#companyForm', '/management-companies'); bindCreateForm('#houseForm', '/houses'); bindCreateForm('#newsForm', '/news');
    document.querySelectorAll('[data-status]').forEach(select => select.onchange = async () => { try { await api(`/requests/${select.dataset.status}/status`, { method: 'PATCH', body: JSON.stringify({ status: select.value }) }); toast('Статус обновлён'); } catch (error) { toast(error.message, true); } });
  } catch (error) { page.innerHTML = `<div class="empty">${escapeHTML(error.message)}</div>`; }
}

function bindCreateForm(selector, path) {
  $(selector).onsubmit = async event => { event.preventDefault(); try { await api(path, { method: 'POST', body: JSON.stringify(Object.fromEntries(new FormData(event.target))) }); toast('Сохранено'); renderManagement(); } catch (error) { toast(error.message, true); } };
}

async function logout(callAPI = true) {
  if (callAPI && state.access) { try { await api('/auth/logout', { method: 'POST', body: JSON.stringify({ refresh_token: state.refresh }) }); } catch (_) {} }
  if (state.notificationController) state.notificationController.abort();
  state.notificationController = null;
  state.access = ''; state.refresh = ''; state.profile = null;
  localStorage.removeItem('access_token'); localStorage.removeItem('refresh_token'); showAuth();
}

$('#authToggle').onclick = () => {
  state.register = !state.register;
  $('#authTitle').textContent = state.register ? 'Регистрация' : 'Вход';
  $('#authSubmit').textContent = state.register ? 'Зарегистрироваться' : 'Войти';
  $('#authToggle').textContent = state.register ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться';
  $('#adminCodeField').classList.toggle('hidden', !state.register);
  $('#password').autocomplete = state.register ? 'new-password' : 'current-password';
};
$('#authForm').onsubmit = async event => {
  event.preventDefault(); const body = { email: $('#email').value.trim(), password: $('#password').value };
  if (state.register && $('#adminCode').value.trim()) body.admin_code = $('#adminCode').value.trim();
  try { setTokens(await api(state.register ? '/auth/register' : '/auth/login', { method: 'POST', body: JSON.stringify(body) })); showApp(); navigate('/profile'); toast(state.register ? 'Аккаунт создан' : 'Добро пожаловать'); }
  catch (error) { toast(error.message, true); }
};
$('#logoutButton').onclick = () => logout();
function deviceToken() {
  let token = localStorage.getItem('notification_device_token');
  if (!token) {
    token = crypto.randomUUID();
    localStorage.setItem('notification_device_token', token);
  }
  return token;
}

async function setNotifications(enabled) {
  if (enabled) {
    if (!('Notification' in window)) throw new Error('Браузер не поддерживает уведомления.');
    const permission = await Notification.requestPermission();
    if (permission !== 'granted') throw new Error('Разрешите уведомления в настройках браузера.');
    await api('/notifications/register', { method: 'POST', body: JSON.stringify({ device_token: deviceToken(), platform: 'web' }) });
    state.notificationsEnabled = true;
    localStorage.setItem('notifications_enabled', 'true');
    startNotificationStream();
  } else {
    if (state.notificationController) state.notificationController.abort();
    await api(`/notifications/unregister?device_token=${encodeURIComponent(deviceToken())}`, { method: 'DELETE' });
    state.notificationsEnabled = false;
    localStorage.setItem('notifications_enabled', 'false');
  }
}

function renderSettings() {
  page.innerHTML = `<div class="grid"><section class="card full"><h3>Уведомления</h3><div class="setting-row"><div><b>Новости и статусы заявок</b><p>Показывать системные уведомления браузера.</p></div><label class="switch"><input id="notificationToggle" type="checkbox" ${state.notificationsEnabled ? 'checked' : ''}><span></span></label></div></section></div>`;
  $('#notificationToggle').onchange = async event => {
    event.target.disabled = true;
    try {
      await setNotifications(event.target.checked);
      toast(event.target.checked ? 'Уведомления включены' : 'Уведомления отключены');
    } catch (error) {
      event.target.checked = state.notificationsEnabled;
      toast(error.message, true);
    } finally {
      event.target.disabled = false;
    }
  };
}

function displayNotification(notification) {
  toast(`${notification.title}: ${notification.body}`);
  if ('Notification' in window && Notification.permission === 'granted' && document.visibilityState !== 'visible') {
    const nativeNotification = new Notification(notification.title, { body: notification.body, tag: `${notification.type}-${notification.created_at}` });
    nativeNotification.onclick = () => { window.focus(); navigate(notification.url || '/profile'); nativeNotification.close(); };
  }
  if (location.pathname === '/news' && notification.type === 'news') renderNews();
  if (location.pathname === '/requests' && notification.type === 'request_status') renderRequests();
}

function parseSSEBlock(block) {
  const lines = block.split('\n');
  const event = lines.find(line => line.startsWith('event:'))?.slice(6).trim();
  const data = lines.filter(line => line.startsWith('data:')).map(line => line.slice(5).trimStart()).join('\n');
  if (event !== 'notification' || !data) return;
  try { displayNotification(JSON.parse(data)); } catch (_) {}
}

async function startNotificationStream() {
  if (!state.access || !state.notificationsEnabled || state.notificationController) return;
  const controller = new AbortController();
  state.notificationController = controller;
  try {
    const response = await fetch(`${API}/notifications/stream`, {
      headers: { Authorization: `Bearer ${state.access}`, Accept: 'text/event-stream' }, signal: controller.signal,
    });
    if (!response.ok || !response.body) throw new Error(`SSE ${response.status}`);
    const reader = response.body.getReader(); const decoder = new TextDecoder(); let buffer = '';
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;
      buffer += decoder.decode(value, { stream: true }).replace(/\r\n/g, '\n');
      let boundary;
      while ((boundary = buffer.indexOf('\n\n')) !== -1) {
        parseSSEBlock(buffer.slice(0, boundary)); buffer = buffer.slice(boundary + 2);
      }
    }
  } catch (error) {
    if (error.name !== 'AbortError') console.warn('Notification stream reconnecting');
  } finally {
    if (state.notificationController === controller) {
      state.notificationController = null;
      if (state.access && state.notificationsEnabled && !controller.signal.aborted) setTimeout(startNotificationStream, 3000);
    }
  }
}
document.addEventListener('click', event => { const link = event.target.closest('[data-route]'); if (link) { event.preventDefault(); navigate(link.getAttribute('href')); } });
window.addEventListener('popstate', () => navigate(location.pathname, false));

if (state.access) { showApp(); navigate(location.pathname, false); } else showAuth();
