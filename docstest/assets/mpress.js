
(() => {
  const doc = document.documentElement;
  const storage = {
    get(key) { try { return localStorage.getItem(key); } catch (_) { return null; } },
    set(key, value) { try { localStorage.setItem(key, value); } catch (_) {} }
  };
  const theme = document.querySelector('#theme');
  const themeModes = ['system', 'dark', 'light'];
  const themeLabels = {system: 'System', dark: 'Dark', light: 'Light'};
  const stored = storage.get('mpress-theme');
  if (themeModes.includes(stored)) doc.dataset.theme = stored;
  if (!themeModes.includes(doc.dataset.theme)) doc.dataset.theme = 'system';
  const syncThemeButton = () => {
    if (!theme) return;
    const current = doc.dataset.theme;
    const next = themeModes[(themeModes.indexOf(current) + 1) % themeModes.length];
    theme.dataset.themeMode = current;
    theme.setAttribute('aria-label', 'Theme: ' + themeLabels[current] + '. Switch to ' + themeLabels[next]);
    theme.title = 'Theme: ' + themeLabels[current];
  };
  syncThemeButton();
  theme?.addEventListener('click', () => {
    const current = themeModes.includes(doc.dataset.theme) ? doc.dataset.theme : 'system';
    const next = themeModes[(themeModes.indexOf(current) + 1) % themeModes.length];
    doc.dataset.theme = next;
    storage.set('mpress-theme', next);
    syncThemeButton();
  });

  const menu = document.querySelector('#menu');
  const sidebar = document.querySelector('.sidebar');
  const mobileNavigation = window.matchMedia('(max-width: 760px)');
  const syncNavigation = () => {
    const closedMobile = mobileNavigation.matches && !document.body.classList.contains('nav-open');
    if (sidebar) {
      sidebar.toggleAttribute('inert', closedMobile);
      sidebar.setAttribute('aria-hidden', String(closedMobile));
    }
  };
  const closeNavigation = restoreFocus => {
    document.body.classList.remove('nav-open');
    menu?.setAttribute('aria-expanded', 'false');
    syncNavigation();
    if (restoreFocus) menu?.focus();
  };
  syncNavigation();
  menu?.addEventListener('click', () => {
    const open = document.body.classList.toggle('nav-open');
    menu.setAttribute('aria-expanded', String(open));
    syncNavigation();
    if (open) sidebar?.focus();
  });
  sidebar?.addEventListener('click', event => {
    if (event.target.closest('a') && window.innerWidth <= 760) {
      closeNavigation(false);
    }
  });
  document.addEventListener('keydown', event => {
    if (event.key === 'Escape' && document.body.classList.contains('nav-open')) closeNavigation(true);
  });
  document.addEventListener('click', event => {
    if (document.body.classList.contains('nav-open') && !event.target.closest('.sidebar') && !event.target.closest('#menu')) closeNavigation(true);
  });
  mobileNavigation.addEventListener?.('change', () => {
    if (!mobileNavigation.matches) closeNavigation(false);
    else syncNavigation();
  });

  const toc = document.querySelector('.toc');
  const tocLinks = toc ? [...toc.querySelectorAll('a[href^="#"]')] : [];
  const tocItems = tocLinks.map(link => ({
    link,
    target: document.getElementById(decodeURIComponent(link.hash.slice(1)))
  })).filter(item => item.target);
  if (tocItems.length) {
    let tocFrame = 0;
    const updateTOC = () => {
      tocFrame = 0;
      const marker = window.scrollY + 112;
      let current = tocItems[0];
      for (const item of tocItems) {
        const top = item.target.getBoundingClientRect().top + window.scrollY;
        if (top <= marker) current = item;
        else break;
      }
      tocItems.forEach(item => {
        const active = item === current;
        item.link.classList.toggle('active', active);
        if (active) item.link.setAttribute('aria-current', 'location');
        else item.link.removeAttribute('aria-current');
      });
      if (toc.clientHeight > 0) {
        const tocRect = toc.getBoundingClientRect();
        const linkRect = current.link.getBoundingClientRect();
        if (linkRect.top < tocRect.top + 12) toc.scrollTop -= tocRect.top + 12 - linkRect.top;
        else if (linkRect.bottom > tocRect.bottom - 12) toc.scrollTop += linkRect.bottom - tocRect.bottom + 12;
      }
    };
    const requestTOCUpdate = () => {
      if (!tocFrame) tocFrame = requestAnimationFrame(updateTOC);
    };
    updateTOC();
    document.addEventListener('scroll', requestTOCUpdate, {passive: true});
    window.addEventListener('resize', requestTOCUpdate, {passive: true});
  }

  const utilityMenus = [...document.querySelectorAll('[data-utility-menu]')];
  const menuButton = menu => document.querySelector('[popovertarget="' + menu.id + '"]');
  const menuIsOpen = menu => menu.matches(':popover-open');
  const positionUtilityMenu = menu => {
    const button = menuButton(menu);
    if (!button) return;
    const anchor = button.getBoundingClientRect();
    const inset = 12;
    const width = menu.offsetWidth;
    const height = menu.offsetHeight;
    const left = Math.max(inset, Math.min(window.innerWidth - width - inset, anchor.right - width));
    const below = anchor.bottom + 8;
    const top = below + height <= window.innerHeight - inset ? below : Math.max(inset, anchor.top - height - 8);
    menu.style.left = left + 'px';
    menu.style.top = top + 'px';
  };
  utilityMenus.forEach(menu => {
    const button = menuButton(menu);
    if (!button) return;
    menu.addEventListener('toggle', event => {
      const open = event.newState ? event.newState === 'open' : menuIsOpen(menu);
      button.setAttribute('aria-expanded', String(open));
      if (open) positionUtilityMenu(menu);
    });
    button.addEventListener('keydown', event => {
      if (event.key !== 'ArrowDown' && event.key !== 'ArrowUp') return;
      event.preventDefault();
      if (!menuIsOpen(menu)) menu.showPopover();
      positionUtilityMenu(menu);
      const items = [...menu.querySelectorAll('a[role="menuitem"]')];
      const target = event.key === 'ArrowUp' ? items[items.length - 1] : (menu.querySelector('[aria-current="true"]') || items[0]);
      target?.focus();
    });
    menu.addEventListener('keydown', event => {
      const items = [...menu.querySelectorAll('a[role="menuitem"]')];
      const current = items.indexOf(document.activeElement);
      let next = current;
      if (event.key === 'ArrowDown') next = (current + 1) % items.length;
      else if (event.key === 'ArrowUp') next = (current - 1 + items.length) % items.length;
      else if (event.key === 'Home') next = 0;
      else if (event.key === 'End') next = items.length - 1;
      else if (event.key === 'Escape') {
        event.preventDefault();
        menu.hidePopover();
        button.focus();
        return;
      } else return;
      event.preventDefault();
      items[next]?.focus();
    });
  });
  const repositionUtilityMenus = () => utilityMenus.forEach(menu => {
    if (menuIsOpen(menu)) positionUtilityMenu(menu);
  });
  window.addEventListener('resize', repositionUtilityMenus, {passive: true});
  document.addEventListener('scroll', repositionUtilityMenus, {passive: true, capture: true});
  document.querySelectorAll('.mpress-tabs, .mpress-preview-tabs').forEach(group => {
    const tabs = [...group.querySelectorAll('[role="tab"]')];
    const panels = [...group.querySelectorAll('[role="tabpanel"]')];
    const select = index => {
      tabs.forEach((item, itemIndex) => item.setAttribute('aria-selected', String(index === itemIndex)));
      tabs.forEach((item, itemIndex) => item.tabIndex = index === itemIndex ? 0 : -1);
      panels.forEach((panel, panelIndex) => panel.hidden = index !== panelIndex);
    };
    tabs.forEach((tab, index) => {
      tab.addEventListener('click', () => select(index));
      tab.addEventListener('keydown', event => {
        let next = index;
        if (event.key === 'ArrowRight') next = (index + 1) % tabs.length;
        else if (event.key === 'ArrowLeft') next = (index - 1 + tabs.length) % tabs.length;
        else if (event.key === 'Home') next = 0;
        else if (event.key === 'End') next = tabs.length - 1;
        else return;
        event.preventDefault();
        select(next);
        tabs[next].focus();
      });
    });
  });

  const blogFilters = [...document.querySelectorAll('[data-blog-filter]')];
  const blogPosts = [...document.querySelectorAll('[data-blog-tags]')];
  blogFilters.forEach(filter => filter.addEventListener('click', () => {
    const selected = filter.dataset.blogFilter;
    blogFilters.forEach(item => {
      const active = item === filter;
      item.classList.toggle('active', active);
      item.setAttribute('aria-pressed', String(active));
    });
    blogPosts.forEach(post => {
      const tags = (post.dataset.blogTags || '').split('|');
      post.hidden = Boolean(selected) && !tags.includes(selected);
    });
  }));

  document.querySelectorAll('.mpress-copy').forEach(button => {
    button.addEventListener('click', async () => {
      const container = button.closest('.mpress-terminal, .mpress-codeframe');
      const source = container?.querySelector('pre');
      const value = source?.dataset.commands || source?.dataset.code || '';
      if (!value) return;
      try {
        await navigator.clipboard.writeText(value);
        const label = button.textContent;
        button.textContent = 'Copied';
        setTimeout(() => button.textContent = label, 1400);
      } catch (_) {
        const field = document.createElement('textarea');
        field.value = value;
        field.style.position = 'fixed';
        field.style.opacity = '0';
        document.body.append(field);
        field.select();
        const copied = document.execCommand('copy');
        field.remove();
        const label = button.textContent;
        button.textContent = copied ? 'Copied' : 'Copy failed';
        setTimeout(() => button.textContent = label, 1400);
      }
    });
  });

  document.querySelectorAll('.mpress-tutorial').forEach(tutorial => {
    const key = 'mpress-' + tutorial.dataset.tutorialId;
    const inputs = [...tutorial.querySelectorAll('input[type="checkbox"]')];
    let saved = [];
    try { saved = JSON.parse(storage.get(key) || '[]'); } catch (_) {}
    const update = () => {
      const done = inputs.filter(input => input.checked).length;
      tutorial.querySelector('.mpress-tutorial-progress-bar')?.style.setProperty('width', (inputs.length ? done / inputs.length * 100 : 0) + '%');
      const count = tutorial.querySelector('.mpress-tutorial-count');
      if (count) count.textContent = done + '/' + inputs.length + ' completed';
      inputs.forEach(input => input.closest('.mpress-tutorial-step')?.classList.toggle('completed', input.checked));
      storage.set(key, JSON.stringify(inputs.filter(input => input.checked).map(input => input.dataset.stepId)));
    };
    inputs.forEach(input => {
      input.checked = saved.includes(input.dataset.stepId);
      input.addEventListener('change', update);
    });
    tutorial.querySelector('.mpress-tutorial-reset')?.addEventListener('click', () => {
      inputs.forEach(input => input.checked = false);
      update();
    });
    update();
  });

  document.querySelectorAll('.mpress-testimonials').forEach(carousel => {
    const items = [...carousel.querySelectorAll('.mpress-testimonial')];
    const dots = [...carousel.querySelectorAll('.mpress-testimonials-dot')];
    let current = 0;
    const show = index => {
      current = index;
      items.forEach((item, i) => {
        item.classList.toggle('active', i === index);
        item.setAttribute('aria-hidden', String(i !== index));
      });
      dots.forEach((dot, i) => i === index ? dot.setAttribute('aria-current', 'true') : dot.removeAttribute('aria-current'));
    };
    dots.forEach((dot, index) => dot.addEventListener('click', () => show(index)));
    const delay = Number(carousel.dataset.autoplay);
    if (items.length > 1 && delay > 0 && !window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      let timer = null;
      const stop = () => { if (timer !== null) clearInterval(timer); timer = null; };
      const start = () => { if (timer === null) timer = setInterval(() => show((current + 1) % items.length), delay); };
      start();
      carousel.addEventListener('mouseenter', stop);
      carousel.addEventListener('mouseleave', start);
      carousel.addEventListener('focusin', stop);
      carousel.addEventListener('focusout', start);
    }
  });

  const audienceBlocks = [...document.querySelectorAll('.mpress-audience[data-audience]')];
  if (audienceBlocks.length) {
    const roles = [...new Set(audienceBlocks.map(block => block.dataset.audience))].sort();
    const selector = document.createElement('div');
    selector.className = 'mpress-audience-selector';
    selector.innerHTML = '<label class="mpress-audience-label">Audience: <select class="mpress-audience-select" aria-label="Select audience role"><option value="all">All</option></select></label>';
    const select = selector.querySelector('select');
    roles.forEach(role => select.add(new Option(role.charAt(0).toUpperCase() + role.slice(1), role)));
    const saved = storage.get('mpress-audience-role');
    select.value = saved && (saved === 'all' || roles.includes(saved)) ? saved : 'all';
    const apply = () => {
      audienceBlocks.forEach(block => block.classList.toggle('mpress-visible', select.value === 'all' || block.dataset.audience === select.value));
      storage.set('mpress-audience-role', select.value);
    };
    select.addEventListener('change', apply);
    audienceBlocks[0].before(selector);
    apply();
  }

  const conditions = [...document.querySelectorAll('.mpress-conditional[data-param]')];
  const params = new URLSearchParams(location.search);
  conditions.forEach(block => {
    const name = block.dataset.param;
    const fromURL = params.get(name);
    if (fromURL) storage.set('mpress-param-' + name, fromURL);
    const value = fromURL || storage.get('mpress-param-' + name) || '';
    block.classList.toggle('mpress-visible', block.dataset.value === value || (!value && block.dataset.default === 'true'));
  });

  const reactiveInputs = [...document.querySelectorAll('[data-reactive]')];
  const reactiveOutputs = [...document.querySelectorAll('.mpress-reactive-computed')];
  if (reactiveInputs.length && reactiveOutputs.length) {
    const store = {};
    const read = input => input.type === 'number' || input.type === 'range' ? Number(input.value) : input.value;
    reactiveInputs.forEach(input => store[input.dataset.reactive] = read(input));
    const evaluate = () => reactiveOutputs.forEach(output => {
      const target = output.querySelector('.mpress-reactive-computed-value');
      try {
        const names = Object.keys(store);
        let result = Function(...names, 'return (' + output.dataset.reactiveExpr + ')')(...names.map(name => store[name]));
        const format = output.dataset.reactiveFormat || '';
        const match = format.match(/%(?:\.(\d+))?([df])/);
        if (match) {
          const value = match[2] === 'd' ? Math.round(Number(result)) : Number(result).toFixed(Number(match[1] || 0));
          result = format.replace(match[0], value);
        }
        target.textContent = String(result);
      } catch (_) { target.textContent = 'Error'; }
    });
    reactiveInputs.forEach(input => input.addEventListener('input', () => {
      store[input.dataset.reactive] = read(input);
      const display = document.querySelector('[data-reactive-display="' + CSS.escape(input.dataset.reactive) + '"]');
      if (display) display.textContent = input.value;
      evaluate();
    }));
    evaluate();
  }

  document.querySelectorAll('.mpress-explained').forEach(panel => {
    const linked = [...panel.querySelectorAll('[data-ref]')];
    const highlight = ref => linked.forEach(item => item.classList.toggle('mpress-explained-active', item.dataset.ref === ref));
    linked.forEach(item => {
      item.addEventListener('mouseenter', () => highlight(item.dataset.ref));
      item.addEventListener('mouseleave', () => highlight(''));
    });
  });

  document.querySelectorAll('.mpress-api-playground').forEach(playground => {
    const url = playground.querySelector('.mpress-api-pg-url-input');
    const server = playground.querySelector('.mpress-api-pg-server-select');
    server?.addEventListener('change', () => url.value = server.value.replace(/\/$/, '') + playground.dataset.path);
    playground.querySelector('.mpress-api-pg-send')?.addEventListener('click', async event => {
      const button = event.currentTarget;
      const responsePanel = playground.querySelector('.mpress-api-pg-response');
      const responseBody = responsePanel.querySelector('code');
      const status = responsePanel.querySelector('.mpress-api-pg-response-status');
      const timing = responsePanel.querySelector('.mpress-api-pg-response-time');
      const headers = {};
      playground.querySelectorAll('.mpress-api-pg-header-row').forEach(row => {
        const fields = row.querySelectorAll('input');
        if (fields[0]?.value) headers[fields[0].value] = fields[1]?.value || '';
      });
      const options = {method: playground.dataset.method, headers};
      const body = playground.querySelector('.mpress-api-pg-request')?.value;
      if (body && !['GET', 'HEAD'].includes(options.method)) options.body = body;
      button.disabled = true;
      const started = performance.now();
      try {
        const response = await fetch(url.value, options);
        const text = await response.text();
        status.textContent = response.status + ' ' + response.statusText;
        timing.textContent = Math.round(performance.now() - started) + ' ms';
        try { responseBody.textContent = JSON.stringify(JSON.parse(text), null, 2); } catch (_) { responseBody.textContent = text; }
      } catch (error) {
        status.textContent = 'Request failed';
        responseBody.textContent = error.message;
      }
      responsePanel.hidden = false;
      button.disabled = false;
    });
  });

  document.querySelectorAll('.mpress-calendar').forEach(calendar => {
    const title = calendar.querySelector('.mpress-calendar-title');
    const grid = calendar.querySelector('.mpress-calendar-grid');
    const draw = (year, month) => {
      calendar.dataset.year = year;
      calendar.dataset.month = month + 1;
      title.textContent = new Intl.DateTimeFormat(undefined, {month: 'long', year: 'numeric'}).format(new Date(year, month, 1));
      [...grid.querySelectorAll('.mpress-calendar-cell')].forEach(cell => cell.remove());
      const first = new Date(year, month, 1);
      const offset = (first.getDay() + 6) % 7;
      const count = new Date(year, month + 1, 0).getDate();
      for (let i = 0; i < offset; i++) {
        const cell = document.createElement('div');
        cell.className = 'mpress-calendar-cell empty';
        grid.append(cell);
      }
      const today = new Date();
      for (let day = 1; day <= count; day++) {
        const cell = document.createElement('div');
        cell.className = 'mpress-calendar-cell';
        if (today.getFullYear() === year && today.getMonth() === month && today.getDate() === day) cell.classList.add('today');
        cell.innerHTML = '<span class="mpress-calendar-date">' + day + '</span>';
        grid.append(cell);
      }
    };
    calendar.querySelectorAll('.mpress-calendar-nav').forEach(button => button.addEventListener('click', () => {
      const date = new Date(Number(calendar.dataset.year), Number(calendar.dataset.month) - 1 + (button.dataset.dir === 'next' ? 1 : -1), 1);
      draw(date.getFullYear(), date.getMonth());
    }));
  });

  const query = document.querySelector('#search');
  if (query) {
    const shortcut = document.querySelector('.search-shortcut');
    const shortcutPlatform = shortcut?.querySelector('.search-shortcut-platform');
    const platform = navigator.userAgentData?.platform || navigator.platform || '';
    const applePlatform = /(Mac|iPhone|iPad|iPod)/i.test(platform);
    if (shortcutPlatform) shortcutPlatform.textContent = applePlatform ? '⌘' : 'Ctrl';
    query.setAttribute('aria-keyshortcuts', applePlatform ? 'Meta+K' : 'Control+K');
    let index = [];
    let indexPromise;
    let active = -1;
    let matches = [];
    let overlay;
    let dialog;
    let overlayInput;
    let resultList;
    let resultCount;
    let preview;
    let debounce;
    const recentKey = 'mpress-recent-searches';
    const escapeHTML = value => String(value || '').replace(/[&<>"']/g, character => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[character]));
    const escapePattern = value => value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    const normalise = value => String(value || '').normalize('NFKD').replace(/[\u0300-\u036f]/g, '').toLowerCase();
    const searchIcon = name => {
      const paths = name === 'clock'
        ? '<circle cx="12" cy="12" r="9"></circle><path d="M12 7v5l3 2"></path>'
        : '<circle cx="11" cy="11" r="8"></circle><path d="m21 21-4.3-4.3"></path>';
      return '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-' + name + '" aria-hidden="true">' + paths + '</svg>';
    };
    const loadIndex = () => {
      if (indexPromise) return indexPromise;
      indexPromise = fetch(document.body.dataset.search)
        .then(response => {
          if (!response.ok) throw new Error('Search index unavailable');
          return response.json();
        })
        .then(value => index = Array.isArray(value) ? value : [])
        .catch(() => index = []);
      return indexPromise;
    };
    const readRecent = () => {
      try {
        const value = JSON.parse(localStorage.getItem(recentKey) || '[]');
        return Array.isArray(value) ? value.filter(item => typeof item === 'string').slice(0, 5) : [];
      } catch (_) {
        return [];
      }
    };
    const saveRecent = value => {
      const clean = value.trim();
      if (!clean) return;
      try {
        localStorage.setItem(recentKey, JSON.stringify([clean, ...readRecent().filter(item => item !== clean)].slice(0, 5)));
      } catch (_) {}
    };
    const levenshtein = (left, right) => {
      if (!left.length) return right.length;
      if (!right.length) return left.length;
      let previous = Array.from({length: left.length + 1}, (_, index) => index);
      for (let row = 1; row <= right.length; row++) {
        const current = [row];
        for (let column = 1; column <= left.length; column++) {
          const cost = left[column - 1] === right[row - 1] ? 0 : 1;
          current[column] = Math.min(current[column - 1] + 1, previous[column] + 1, previous[column - 1] + cost);
        }
        previous = current;
      }
      return previous[left.length];
    };
    const fuzzyMatch = (word, candidate) => {
      if (candidate.includes(word)) return true;
      if (word.length < 4) return false;
      return levenshtein(word, candidate) <= (word.length >= 6 ? 2 : 1);
    };
    const scoreItem = (item, words, fullQuery) => {
      const title = normalise(item.title);
      const description = normalise(item.description);
      const text = normalise(item.text);
      const headings = (item.headings || []).map(heading => normalise(heading.text || heading));
      const tags = (item.tags || []).map(normalise);
      let score = 0;
      if (title === fullQuery) score += 160;
      else if (title.startsWith(fullQuery)) score += 110;
      else if (title.includes(fullQuery)) score += 75;
      if (headings.some(heading => heading.includes(fullQuery))) score += 45;
      for (const word of words) {
        let wordScore = 0;
        if (title.includes(word)) wordScore = Math.max(wordScore, 34 + (title.startsWith(word) ? 12 : 0));
        if (headings.some(heading => heading.includes(word))) wordScore = Math.max(wordScore, 24);
        if (tags.some(tag => tag.includes(word))) wordScore = Math.max(wordScore, 20);
        if (description.includes(word)) wordScore = Math.max(wordScore, 14);
        if (text.includes(word)) wordScore = Math.max(wordScore, 5);
        if (!wordScore) {
          const titleWords = title.split(/\s+/);
          const headingWords = headings.flatMap(heading => heading.split(/\s+/));
          if (titleWords.some(candidate => fuzzyMatch(word, candidate))) wordScore = 14;
          else if (headingWords.some(candidate => fuzzyMatch(word, candidate))) wordScore = 8;
        }
        if (!wordScore) return -1;
        score += wordScore;
      }
      return score;
    };
    const highlight = (value, words) => {
      let safe = escapeHTML(value);
      const unique = [...new Set(words)].sort((left, right) => right.length - left.length);
      if (!unique.length) return safe;
      const pattern = new RegExp('(' + unique.map(escapePattern).join('|') + ')', 'gi');
      return safe.replace(pattern, '<mark>$1</mark>');
    };
    const snippetFor = (item, words, limit = 190) => {
      const source = String(item.description || item.text || '').replace(/\s+/g, ' ').trim();
      if (!source) return '';
      const lower = normalise(source);
      let start = words.reduce((best, word) => {
        const position = lower.indexOf(word);
        return position >= 0 && (best < 0 || position < best) ? position : best;
      }, -1);
      if (start < 0) start = 0;
      start = Math.max(0, start - 55);
      let excerpt = source.slice(start, start + limit).trim();
      if (start > 0) excerpt = '…' + excerpt;
      if (start + limit < source.length) excerpt += '…';
      return highlight(excerpt, words);
    };
    const ensureOverlay = () => {
      if (overlay) return;
      overlay = document.createElement('div');
      overlay.className = 'mpress-search-overlay';
      overlay.innerHTML = '<section class="mpress-search-dialog" id="mpress-search-dialog" role="dialog" aria-modal="true" aria-label="Search documentation">' +
        '<div class="mpress-search-query">' + searchIcon('search') + '<input type="search" autocomplete="off" spellcheck="false" placeholder="Search documentation" aria-label="Search documentation" role="combobox" aria-expanded="false" aria-controls="mpress-search-listbox" aria-autocomplete="list"><button class="mpress-search-close" type="button" aria-label="Close search">ESC</button></div>' +
        '<div class="mpress-search-meta"><span class="mpress-search-count" aria-live="polite">Type to search</span><span class="mpress-search-hint"><kbd>↑</kbd><kbd>↓</kbd> navigate · <kbd>Enter</kbd> open</span></div>' +
        '<div class="mpress-search-list" id="mpress-search-listbox" role="listbox" aria-label="Search results"></div><aside class="mpress-search-preview" aria-label="Result preview"></aside></section>';
      document.body.append(overlay);
      dialog = overlay.querySelector('.mpress-search-dialog');
      overlayInput = overlay.querySelector('.mpress-search-query input');
      resultList = overlay.querySelector('.mpress-search-list');
      resultCount = overlay.querySelector('.mpress-search-count');
      preview = overlay.querySelector('.mpress-search-preview');
      overlay.addEventListener('click', event => {
        if (event.target === overlay) closeSearch();
      });
      overlay.querySelector('.mpress-search-close').addEventListener('click', closeSearch);
      overlayInput.addEventListener('input', () => {
        clearTimeout(debounce);
        debounce = setTimeout(runSearch, 45);
      });
      overlayInput.addEventListener('keydown', event => {
        if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
          event.preventDefault();
          setActive(active + (event.key === 'ArrowDown' ? 1 : -1));
        } else if (event.key === 'Enter' && matches.length) {
          event.preventDefault();
          const target = resultList.querySelector('#mpress-search-option-' + (active < 0 ? 0 : active));
          target?.click();
        }
      });
    };
    const renderPreview = (item, words) => {
      if (!item) {
        preview.replaceChildren();
        return;
      }
      let html = '<h2>' + highlight(item.title, words) + '</h2>';
      if (item.description) html += '<p>' + highlight(item.description, words) + '</p>';
      const headings = (item.headings || []).filter(heading => heading.text).slice(0, 9);
      if (headings.length) {
        html += '<ul class="mpress-search-preview-headings">';
        headings.forEach(heading => {
          html += '<li class="level-' + Number(heading.level || 2) + '"><a href="' + escapeHTML(item.url) + '#' + encodeURIComponent(heading.id || '') + '">' + highlight(heading.text, words) + '</a></li>';
        });
        html += '</ul>';
      }
      const body = snippetFor(item, words, 620);
      if (body) html += '<div class="mpress-search-preview-body">' + body + '</div>';
      preview.innerHTML = html;
    };
    const resultLinks = () => [...resultList.querySelectorAll('a.mpress-search-item')];
    const setActive = value => {
      const links = resultLinks();
      active = links.length ? Math.max(0, Math.min(value, links.length - 1)) : -1;
      links.forEach((link, position) => {
        const selected = position === active;
        link.classList.toggle('active', selected);
        link.setAttribute('aria-selected', String(selected));
      });
      if (active >= 0) {
        overlayInput.setAttribute('aria-activedescendant', links[active].id);
        links[active].scrollIntoView({block: 'nearest'});
        renderPreview(matches[active]?.item, normalise(overlayInput.value).split(/\s+/).filter(Boolean));
      } else {
        overlayInput.removeAttribute('aria-activedescendant');
        renderPreview(null, []);
      }
    };
    const renderRecent = () => {
      ensureOverlay();
      const recent = readRecent();
      matches = [];
      active = -1;
      dialog.classList.remove('has-results');
      overlayInput.setAttribute('aria-expanded', 'false');
      overlayInput.removeAttribute('aria-activedescendant');
      resultCount.textContent = recent.length ? 'Recent searches' : 'Type to search';
      preview.replaceChildren();
      if (!recent.length) {
        resultList.innerHTML = '<div class="mpress-search-empty"><strong>Search this documentation</strong><span>Find pages, headings, concepts, and code terms.</span></div>';
        return;
      }
      resultList.innerHTML = '<div class="mpress-search-recent-heading"><span>Recent</span><button class="mpress-search-clear" type="button">Clear</button></div><div class="mpress-search-recent">' + recent.map((value, position) => '<button class="mpress-search-item" type="button" data-recent="' + position + '">' + searchIcon('clock') + '<span>' + escapeHTML(value) + '</span></button>').join('') + '</div>';
      resultList.querySelector('.mpress-search-clear').addEventListener('click', () => {
        try { localStorage.removeItem(recentKey); } catch (_) {}
        renderRecent();
      });
      resultList.querySelectorAll('[data-recent]').forEach(button => button.addEventListener('click', () => {
        overlayInput.value = recent[Number(button.dataset.recent)] || '';
        runSearch();
      }));
    };
    const renderResults = (queryValue, words) => {
      const resultTotal = matches.length;
      resultCount.textContent = resultTotal ? resultTotal + ' result' + (resultTotal === 1 ? '' : 's') : 'No results';
      dialog.classList.toggle('has-results', resultTotal > 0);
      overlayInput.setAttribute('aria-expanded', String(resultTotal > 0));
      if (!resultTotal) {
        active = -1;
        overlayInput.removeAttribute('aria-activedescendant');
        preview.replaceChildren();
        resultList.innerHTML = '<div class="mpress-search-empty"><strong>No pages found for “' + escapeHTML(queryValue) + '”</strong><span>Try fewer words or check the spelling.</span></div>';
        return;
      }
      resultList.innerHTML = matches.map((match, position) => {
        const item = match.item;
        return '<a class="mpress-search-item" id="mpress-search-option-' + position + '" href="' + escapeHTML(item.url) + '" role="option" aria-selected="false"><div class="mpress-search-item-title"><span>' + highlight(item.title, words) + '</span><span class="mpress-search-item-path">' + escapeHTML(item.url) + '</span></div><div class="mpress-search-item-snippet">' + snippetFor(item, words) + '</div></a>';
      }).join('');
      resultLinks().forEach((link, position) => {
        link.addEventListener('mouseenter', () => setActive(position));
        link.addEventListener('focus', () => setActive(position));
        link.addEventListener('click', () => saveRecent(queryValue));
      });
      setActive(0);
    };
    const runSearch = async () => {
      ensureOverlay();
      const queryValue = overlayInput.value.trim();
      if (!queryValue) {
        renderRecent();
        return;
      }
      await loadIndex();
      if (overlayInput.value.trim() !== queryValue) return;
      const fullQuery = normalise(queryValue);
      const words = fullQuery.split(/\s+/).filter(Boolean);
      matches = index.map(item => ({item, score: scoreItem(item, words, fullQuery)}))
        .filter(match => match.score >= 0)
        .sort((left, right) => right.score - left.score || left.item.title.localeCompare(right.item.title))
        .slice(0, 12);
      renderResults(queryValue, words);
    };
    const openSearch = () => {
      ensureOverlay();
      overlay.classList.add('open');
      document.body.classList.add('mpress-search-open');
      query.setAttribute('aria-expanded', 'true');
      renderRecent();
      void loadIndex();
      requestAnimationFrame(() => overlayInput.focus());
    };
    function closeSearch() {
      if (!overlay?.classList.contains('open')) return;
      overlay.classList.remove('open');
      document.body.classList.remove('mpress-search-open');
      query.setAttribute('aria-expanded', 'false');
      query.removeAttribute('aria-activedescendant');
      overlayInput.value = '';
      matches = [];
      active = -1;
      query.blur();
    }
    query.setAttribute('aria-expanded', 'false');
    query.addEventListener('focus', openSearch);
    query.addEventListener('click', openSearch);
    document.addEventListener('keydown', event => {
      if ((applePlatform ? event.metaKey : event.ctrlKey) && event.key.toLowerCase() === 'k') {
        event.preventDefault();
        if (overlay?.classList.contains('open')) overlayInput.focus();
        else openSearch();
      } else if (event.key === 'Escape') {
        closeSearch();
      }
    });
  }

})();
