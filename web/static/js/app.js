// API Configuration
const API_BASE_URL = window.location.origin; // Use same origin as frontend

// State Management
let authToken = localStorage.getItem('authToken');
let currentUser = JSON.parse(localStorage.getItem('currentUser') || 'null');

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    updateAuthUI();
    // Set minimum date for expiration input
    const expiresAtInput = document.getElementById('expiresAtInput');
    if (expiresAtInput) {
        const now = new Date();
        now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
        expiresAtInput.min = now.toISOString().slice(0, 16);
    }
    // Initialize theme
    initTheme();
});

// Dark Mode Functions
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
    updateThemeIcon(savedTheme);
}

function toggleDarkMode() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    updateThemeIcon(newTheme);
}

function updateThemeIcon(theme) {
    const toggle = document.getElementById('themeToggle');
    if (toggle) {
        toggle.textContent = theme === 'dark' ? '☀️' : '🌙';
    }
}

// Auth Functions
function showLogin() {
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('authModal').classList.add('show');
}

function showRegister() {
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('registerForm').style.display = 'block';
    document.getElementById('authModal').classList.add('show');
}

function closeAuthModal() {
    document.getElementById('authModal').classList.remove('show');
}

async function handleRegister(event) {
    event.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch(`${API_BASE_URL}/api/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        });

        const data = await response.json();

        if (response.ok) {
            authToken = data.token;
            currentUser = data.user;
            localStorage.setItem('authToken', authToken);
            localStorage.setItem('currentUser', JSON.stringify(currentUser));
            updateAuthUI();
            closeAuthModal();
            showNotification('Registration successful!', 'success');
        } else {
            showNotification(data.error || 'Registration failed', 'error');
        }
    } catch (error) {
        showNotification('Network error. Please try again.', 'error');
    }
}

async function handleLogin(event) {
    event.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();

        if (response.ok) {
            authToken = data.token;
            currentUser = data.user;
            localStorage.setItem('authToken', authToken);
            localStorage.setItem('currentUser', JSON.stringify(currentUser));
            updateAuthUI();
            closeAuthModal();
            showNotification('Login successful!', 'success');
        } else {
            showNotification(data.error || 'Login failed', 'error');
        }
    } catch (error) {
        showNotification('Network error. Please try again.', 'error');
    }
}

function logout() {
    authToken = null;
    currentUser = null;
    localStorage.removeItem('authToken');
    localStorage.removeItem('currentUser');
    updateAuthUI();
    showNotification('Logged out successfully', 'success');
}

function updateAuthUI() {
    const authButtons = document.getElementById('authButtons');
    const userInfo = document.getElementById('userInfo');
    const usernameSpan = document.getElementById('username');
    const myUrlsSection = document.getElementById('myUrlsSection');

    if (authToken && currentUser) {
        authButtons.style.display = 'none';
        userInfo.style.display = 'flex';
        usernameSpan.textContent = `👤 ${currentUser.username}`;
        myUrlsSection.style.display = 'block';
        loadMyURLs();
    } else {
        authButtons.style.display = 'flex';
        userInfo.style.display = 'none';
        myUrlsSection.style.display = 'none';
    }
}

// URL Shortening
async function handleShorten(event) {
    event.preventDefault();
    const urlInput = document.getElementById('urlInput');
    const customCodeInput = document.getElementById('customCodeInput');
    const expiresAtInput = document.getElementById('expiresAtInput');
    const url = urlInput.value;
    const customCode = customCodeInput.value.trim();
    const expiresAt = expiresAtInput.value;
    const resultDiv = document.getElementById('result');

    resultDiv.innerHTML = '<div class="loading"></div> Loading...';
    resultDiv.className = 'result show';

    try {
        const headers = {
            'Content-Type': 'application/json'
        };

        if (authToken) {
            headers['Authorization'] = `Bearer ${authToken}`;
        }

        const body = { url };
        if (customCode) {
            body.custom_code = customCode;
        }
        const expirationCheckbox = document.getElementById('enableExpiration');
        if (expirationCheckbox && expirationCheckbox.checked && expiresAt) {
            body.expires_at = new Date(expiresAt).toISOString();
        }

        const response = await fetch(`${API_BASE_URL}/api/shorten`, {
            method: 'POST',
            headers: headers,
            body: JSON.stringify(body)
        });

        const data = await response.json();

        if (response.ok) {
            const shortUrl = `${API_BASE_URL}/${data.code}`;
            const qrUrl = `${API_BASE_URL}/api/qr/${data.code}?size=300`;
            
            // Store QR code URL globally for copy/download
            window.currentQRCodeUrl = qrUrl;
            window.currentQRCodeData = data.code;
            
            resultDiv.innerHTML = `
                <h3>✅ URL Shortened Successfully!</h3>
                <div class="result-content">
                    <div class="result-item">
                        <span class="result-label">Original URL:</span>
                        <a href="${data.original_url}" target="_blank" class="result-link">${data.original_url}</a>
                    </div>
                    <div class="result-item">
                        <span class="result-label">Short URL:</span>
                        <a href="${shortUrl}" target="_blank" class="result-link short-url">${shortUrl}</a>
                    </div>
                    ${data.expires_at ? `<div class="result-item"><span class="result-label">Expires:</span><span>${new Date(data.expires_at).toLocaleString()}</span></div>` : ''}
                </div>
                <div class="qr-section">
                    <div class="qr-preview">
                        <img src="${qrUrl}" alt="QR Code" id="qrCodeImage" class="qr-image" onclick="showQRModal()" oncontextmenu="return false;">
                        <p class="qr-hint">Click QR code to view & copy • Right-click to save</p>
                    </div>
                    <div class="action-buttons">
                        <button class="btn btn-primary" onclick="copyToClipboard('${shortUrl}')">📋 Copy Short URL</button>
                        <button class="btn btn-secondary" onclick="showQRModal()">🔲 View QR Code</button>
                        <button class="btn btn-secondary" onclick="downloadQRCode()" style="font-size: 0.9rem;">💾 Download QR</button>
                    </div>
                </div>
            `;
            resultDiv.className = 'result show success';
            urlInput.value = '';
            customCodeInput.value = '';
            expiresAtInput.value = '';
            document.getElementById('enableExpiration').checked = false;
            toggleExpiration();
            
            // Refresh My URLs if logged in
            if (authToken) {
                loadMyURLs();
            }
        } else {
            resultDiv.innerHTML = `<h3>❌ Error</h3><p>${data.error || 'Failed to shorten URL'}</p>`;
            resultDiv.className = 'result show error';
        }
    } catch (error) {
        resultDiv.innerHTML = `<h3>❌ Error</h3><p>Network error. Please try again.</p>`;
        resultDiv.className = 'result show error';
    }
}

// Store current stats data for export
let currentStatsData = null;
let currentStatsCode = null;

// Stats Lookup
async function handleStats(event) {
    event.preventDefault();
    const codeInput = document.getElementById('codeInput');
    let code = codeInput.value.trim();
    
    // Extract code from full URL if user pasted full URL
    if (code.includes('http://') || code.includes('https://')) {
        // Extract code from URL (e.g., http://localhost:8080/WO8iDmzy -> WO8iDmzy)
        const urlParts = code.split('/');
        code = urlParts[urlParts.length - 1];
        // Remove any query parameters or fragments
        code = code.split('?')[0].split('#')[0];
    }
    
    // Remove whitespace
    code = code.trim();
    
    if (!code) {
        const statsDiv = document.getElementById('statsResult');
        statsDiv.innerHTML = `<div class="result error"><h3>❌ Error</h3><p>Please enter a valid short code</p></div>`;
        statsDiv.className = 'stats-result show';
        return;
    }
    
    const statsDiv = document.getElementById('statsResult');
    const exportButtons = document.getElementById('exportButtons');
    statsDiv.innerHTML = '<div class="loading"></div> Loading stats...';
    statsDiv.className = 'stats-result show';
    exportButtons.style.display = 'none';

    try {
        // Try enhanced stats first
        const response = await fetch(`${API_BASE_URL}/api/stats/${code}/enhanced`);
        const data = await response.json();

        if (response.ok) {
            currentStatsData = data;
            currentStatsCode = code;
            displayEnhancedStats(data, statsDiv);
            exportButtons.style.display = 'block';
        } else {
            // Fallback to basic stats
            const basicResponse = await fetch(`${API_BASE_URL}/api/stats/${code}`);
            const basicData = await basicResponse.json();
            
            if (basicResponse.ok) {
                currentStatsData = basicData;
                currentStatsCode = code;
                displayBasicStats(basicData, statsDiv);
                exportButtons.style.display = 'block';
            } else {
                currentStatsData = null;
                currentStatsCode = null;
                statsDiv.innerHTML = `<div class="result error"><h3>❌ Error</h3><p>${basicData.error || 'Stats not found'}</p></div>`;
                exportButtons.style.display = 'none';
            }
        }
    } catch (error) {
        console.error('Stats error:', error);
        currentStatsData = null;
        currentStatsCode = null;
        statsDiv.innerHTML = `<div class="result error"><h3>❌ Error</h3><p>Network error. Please try again.</p></div>`;
        exportButtons.style.display = 'none';
    }
}

function displayBasicStats(data, container) {
    container.innerHTML = `
        <div class="stat-card">
            <h3>📊 Basic Statistics</h3>
            <div class="stat-row">
                <span class="stat-label">Original URL:</span>
                <span class="stat-value"><a href="${data.original_url}" target="_blank">${data.original_url}</a></span>
            </div>
            <div class="stat-row">
                <span class="stat-label">Total Clicks:</span>
                <span class="stat-value">${data.total_clicks}</span>
            </div>
            <div class="stat-row">
                <span class="stat-label">Unique Visitors:</span>
                <span class="stat-value">${data.unique_ips}</span>
            </div>
            <div class="stat-row">
                <span class="stat-label">Created:</span>
                <span class="stat-value">${new Date(data.created_at).toLocaleDateString()}</span>
            </div>
        </div>
    `;
}

function displayEnhancedStats(data, container) {
    // Check if data is valid
    if (!data || !data.code) {
        container.innerHTML = `<div class="result error"><h3>❌ Error</h3><p>Invalid stats data received</p></div>`;
        return;
    }
    
    let clicksByDayHTML = '';
    if (data.clicks_by_day && Object.keys(data.clicks_by_day).length > 0) {
        const maxClicks = Math.max(...Object.values(data.clicks_by_day));
        clicksByDayHTML = '<div class="chart-container"><h4>📈 Clicks by Day (Last 30 Days)</h4><div class="bar-chart">';
        Object.entries(data.clicks_by_day)
            .sort((a, b) => new Date(a[0]) - new Date(b[0]))
            .forEach(([day, count]) => {
                const percentage = maxClicks > 0 ? (count / maxClicks) * 100 : 0;
                clicksByDayHTML += `
                    <div class="bar-item">
                        <span class="bar-label">${new Date(day).toLocaleDateString()}</span>
                        <div class="bar" style="width: ${percentage}%">${count}</div>
                    </div>
                `;
            });
        clicksByDayHTML += '</div></div>';
    }

    let referrersHTML = '';
    if (data.top_referrers && data.top_referrers.length > 0) {
        referrersHTML = '<div class="stat-card"><h4>🔗 Top Referrers</h4>';
        data.top_referrers.forEach(ref => {
            referrersHTML += `
                <div class="stat-row">
                    <span class="stat-label">${ref.referrer || 'Direct'}</span>
                    <span class="stat-value">${ref.count}</span>
                </div>
            `;
        });
        referrersHTML += '</div>';
    }

    let userAgentsHTML = '';
    if (data.user_agents && Object.keys(data.user_agents).length > 0) {
        const maxUA = Math.max(...Object.values(data.user_agents));
        userAgentsHTML = '<div class="chart-container"><h4>🌐 User Agents</h4><div class="bar-chart">';
        Object.entries(data.user_agents)
            .sort((a, b) => b[1] - a[1])
            .slice(0, 10)
            .forEach(([ua, count]) => {
                const percentage = maxUA > 0 ? (count / maxUA) * 100 : 0;
                userAgentsHTML += `
                    <div class="bar-item">
                        <span class="bar-label">${ua}</span>
                        <div class="bar" style="width: ${percentage}%">${count}</div>
                    </div>
                `;
            });
        userAgentsHTML += '</div></div>';
    }

    let countriesHTML = '';
    if (data.countries && Object.keys(data.countries).length > 0) {
        const maxCountry = Math.max(...Object.values(data.countries));
        const totalCountryClicks = Object.values(data.countries).reduce((a, b) => a + b, 0);
        const hasLocalOnly = Object.keys(data.countries).every(c => c.includes('Local'));
        
        countriesHTML = '<div class="chart-container"><h4>🌍 Geographic Distribution</h4>';
        if (hasLocalOnly) {
            countriesHTML += '<p style="color: var(--text-light); font-size: 0.9em; margin-bottom: 15px;">💡 <em>Showing "Local" because you\'re testing from localhost. In production with real visitors, you\'ll see actual countries!</em></p>';
        }
        countriesHTML += '<div class="bar-chart">';
        Object.entries(data.countries)
            .sort((a, b) => b[1] - a[1])
            .slice(0, 15)
            .forEach(([country, count]) => {
                const percentage = maxCountry > 0 ? (count / maxCountry) * 100 : 0;
                const countryPercentage = totalCountryClicks > 0 ? ((count / totalCountryClicks) * 100).toFixed(1) : 0;
                countriesHTML += `
                    <div class="bar-item">
                        <span class="bar-label">${country} (${countryPercentage}%)</span>
                        <div class="bar" style="width: ${percentage}%">${count}</div>
                    </div>
                `;
            });
        countriesHTML += '</div></div>';
    } else {
        // Show message if no country data available
        countriesHTML = '<div class="stat-card"><h4>🌍 Geographic Distribution</h4><p style="color: var(--text-light); font-style: italic;">No geographic data available yet. Country tracking starts for new clicks after this feature was added.</p></div>';
    }

    container.innerHTML = `
        <div class="stat-card">
            <h3>📊 Enhanced Statistics</h3>
            <div class="stat-row">
                <span class="stat-label">Original URL:</span>
                <span class="stat-value"><a href="${data.original_url}" target="_blank">${data.original_url}</a></span>
            </div>
            <div class="stat-row">
                <span class="stat-label">Total Clicks:</span>
                <span class="stat-value">${data.total_clicks || 0}</span>
            </div>
            <div class="stat-row">
                <span class="stat-label">Unique Visitors:</span>
                <span class="stat-value">${data.unique_ips || 0}</span>
            </div>
            <div class="stat-row">
                <span class="stat-label">Created:</span>
                <span class="stat-value">${new Date(data.created_at).toLocaleDateString()}</span>
            </div>
        </div>
        ${clicksByDayHTML}
        ${referrersHTML}
        ${userAgentsHTML}
        ${countriesHTML}
    `;
}

// My URLs Management
async function loadMyURLs() {
    if (!authToken) return;
    
    const myUrlsList = document.getElementById('myUrlsList');
    myUrlsList.innerHTML = '<div class="loading"></div> Loading your URLs...';

    try {
        const response = await fetch(`${API_BASE_URL}/api/my-urls`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        const data = await response.json();

        if (response.ok) {
            if (data.urls && data.urls.length > 0) {
                let html = `<p style="margin-bottom: 15px; color: var(--text-light);">You have ${data.count} shortened URL(s)</p>`;
                html += '<div style="display: flex; flex-direction: column; gap: 15px;">';
                
                data.urls.forEach(url => {
                    const shortUrl = `${API_BASE_URL}/${url.code}`;
                    const qrUrl = `${API_BASE_URL}/api/qr/${url.code}`;
                    const createdDate = new Date(url.created_at).toLocaleDateString();
                    
                    html += `
                        <div style="border: 1px solid var(--border-color); border-radius: 8px; padding: 15px; background: var(--bg-color);">
                            <div style="display: flex; justify-content: space-between; align-items: start; flex-wrap: wrap; gap: 10px;">
                                <div style="flex: 1; min-width: 200px;">
                                    <p style="margin: 0 0 5px 0; font-weight: 600; color: var(--text-color);">
                                        <a href="${shortUrl}" target="_blank">${shortUrl}</a>
                                    </p>
                                    <p style="margin: 0 0 5px 0; font-size: 0.9rem; color: var(--text-light); word-break: break-all;">
                                        → ${url.original_url}
                                    </p>
                                    <p style="margin: 0; font-size: 0.85rem; color: var(--text-light);">
                                        Created: ${createdDate}
                                    </p>
                                </div>
                                <div style="display: flex; gap: 8px; flex-wrap: wrap;">
                                    <img src="${qrUrl}" alt="QR Code" style="width: 60px; height: 60px; border: 1px solid var(--border-color); border-radius: 4px; padding: 5px; background: white;">
                                    <div style="display: flex; flex-direction: column; gap: 5px;">
                                        <button class="btn btn-secondary" style="padding: 8px 12px; font-size: 0.85rem;" onclick="copyToClipboard('${shortUrl}')">Copy</button>
                                        <button class="btn btn-secondary" style="padding: 8px 12px; font-size: 0.85rem;" onclick="viewStats('${url.code}')">Stats</button>
                                        <button class="btn btn-secondary" style="padding: 8px 12px; font-size: 0.85rem; background: var(--error-color);" onclick="deleteURL('${url.code}')">Delete</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    `;
                });
                
                html += '</div>';
                myUrlsList.innerHTML = html;
            } else {
                myUrlsList.innerHTML = '<p style="color: var(--text-light); text-align: center; padding: 20px;">No URLs yet. Create your first shortened URL above!</p>';
            }
        } else {
            myUrlsList.innerHTML = `<div class="result error"><p>Error loading URLs: ${data.error || 'Unknown error'}</p></div>`;
        }
    } catch (error) {
        myUrlsList.innerHTML = '<div class="result error"><p>Network error. Please try again.</p></div>';
    }
}

async function deleteURL(code) {
    if (!confirm(`Are you sure you want to delete ${code}?`)) return;
    
    if (!authToken) return;

    try {
        const response = await fetch(`${API_BASE_URL}/api/urls/${code}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        const data = await response.json();

        if (response.ok) {
            showNotification('URL deleted successfully', 'success');
            loadMyURLs();
        } else {
            showNotification(data.error || 'Failed to delete URL', 'error');
        }
    } catch (error) {
        showNotification('Network error. Please try again.', 'error');
    }
}

function viewStats(code) {
    document.getElementById('codeInput').value = code;
    document.getElementById('statsForm').dispatchEvent(new Event('submit'));
    // Scroll to stats section
    document.getElementById('statsResult').scrollIntoView({ behavior: 'smooth', block: 'start' });
}

// Export Analytics Functions
function exportStats(format) {
    if (!currentStatsData || !currentStatsCode) {
        showNotification('No stats data to export', 'error');
        return;
    }
    
    let content, filename, mimeType;
    
    if (format === 'csv') {
        content = convertToCSV(currentStatsData);
        filename = `stats-${currentStatsCode}.csv`;
        mimeType = 'text/csv';
    } else {
        content = JSON.stringify(currentStatsData, null, 2);
        filename = `stats-${currentStatsCode}.json`;
        mimeType = 'application/json';
    }
    
    // Create download
    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    
    showNotification(`${format.toUpperCase()} exported successfully!`, 'success');
}

function convertToCSV(data) {
    let csv = 'Metric,Value\n';
    csv += `Code,${data.code}\n`;
    csv += `Original URL,${data.original_url}\n`;
    csv += `Created At,${data.created_at}\n`;
    csv += `Total Clicks,${data.total_clicks || 0}\n`;
    csv += `Unique IPs,${data.unique_ips || 0}\n`;
    
    if (data.clicks_by_day) {
        csv += '\nClicks by Day\n';
        csv += 'Date,Clicks\n';
        Object.entries(data.clicks_by_day)
            .sort((a, b) => new Date(a[0]) - new Date(b[0]))
            .forEach(([day, count]) => {
                csv += `${day},${count}\n`;
            });
    }
    
    if (data.top_referrers && data.top_referrers.length > 0) {
        csv += '\nTop Referrers\n';
        csv += 'Referrer,Count\n';
        data.top_referrers.forEach(ref => {
            csv += `${ref.referrer || 'Direct'},${ref.count}\n`;
        });
    }
    
    if (data.user_agents) {
        csv += '\nUser Agents\n';
        csv += 'User Agent,Count\n';
        Object.entries(data.user_agents)
            .sort((a, b) => b[1] - a[1])
            .forEach(([ua, count]) => {
                csv += `${ua},${count}\n`;
            });
    }
    
    if (data.countries) {
        csv += '\nCountries\n';
        csv += 'Country,Count\n';
        Object.entries(data.countries)
            .sort((a, b) => b[1] - a[1])
            .forEach(([country, count]) => {
                csv += `${country},${count}\n`;
            });
    }
    
    return csv;
}

// Utility Functions
function toggleExpiration() {
    const checkbox = document.getElementById('enableExpiration');
    const input = document.getElementById('expiresAtInput');
    
    if (checkbox.checked) {
        input.style.display = 'block';
        // Set minimum date to today
        const now = new Date();
        now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
        input.min = now.toISOString().slice(0, 16);
    } else {
        input.style.display = 'none';
        input.value = '';
    }
}

function showQRModal() {
    const modal = document.getElementById('qrModal');
    const container = document.getElementById('qrCodeContainer');
    
    if (window.currentQRCodeUrl) {
        // Create image element with onload handler
        const img = document.createElement('img');
        img.src = window.currentQRCodeUrl;
        img.alt = 'QR Code';
        img.style.maxWidth = '100%';
        img.style.borderRadius = '8px';
        img.id = 'modalQRImage';
        
        // Wait for image to load before showing modal
        img.onload = () => {
            container.innerHTML = '';
            container.appendChild(img);
            modal.classList.add('show');
        };
        
        img.onerror = () => {
            container.innerHTML = '<p style="color: var(--error-color);">Failed to load QR code</p>';
            modal.classList.add('show');
        };
        
        // Preload image
        container.appendChild(img);
        
        // If already loaded, show immediately
        if (img.complete) {
            modal.classList.add('show');
        }
    }
}

function closeQRModal() {
    document.getElementById('qrModal').classList.remove('show');
}

async function copyQRCode() {
    if (!window.currentQRCodeUrl) {
        showNotification('No QR code available', 'error');
        return;
    }
    
    try {
        // Try to get image from modal first, then from result
        let img = document.getElementById('modalQRImage');
        if (!img) {
            img = document.getElementById('qrCodeImage');
        }
        
        if (!img) {
            // Fetch image if not in DOM
            const response = await fetch(window.currentQRCodeUrl);
            if (!response.ok) throw new Error('Failed to fetch QR code');
            
            const blob = await response.blob();
            
            // Check if ClipboardItem is supported
            if (navigator.clipboard && window.ClipboardItem) {
                await navigator.clipboard.write([
                    new ClipboardItem({ [blob.type]: blob })
                ]);
                showNotification('QR Code copied to clipboard!', 'success');
                return;
            } else {
                throw new Error('ClipboardItem not supported');
            }
        }
        
        // Use canvas method (more reliable)
        if (!img.complete) {
            showNotification('Please wait for QR code to load...', 'error');
            return;
        }
        
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        canvas.width = img.naturalWidth || img.width;
        canvas.height = img.naturalHeight || img.height;
        
        // Draw image to canvas
        ctx.drawImage(img, 0, 0);
        
        // Convert to blob and copy
        canvas.toBlob(async (blob) => {
            if (!blob) {
                showNotification('Failed to process QR code image', 'error');
                return;
            }
            
            try {
                // Try modern ClipboardItem API
                if (navigator.clipboard && window.ClipboardItem) {
                    await navigator.clipboard.write([
                        new ClipboardItem({ [blob.type]: blob })
                    ]);
                    showNotification('QR Code copied to clipboard!', 'success');
                } else {
                    // Fallback: copy as data URL (some browsers)
                    const dataURL = canvas.toDataURL('image/png');
                    await navigator.clipboard.writeText(dataURL);
                    showNotification('QR Code data copied! Paste in image editor.', 'success');
                }
            } catch (err) {
                console.error('Copy error:', err);
                // Last resort: show data URL for manual copy
                const dataURL = canvas.toDataURL('image/png');
                prompt('Copy this image data (right-click to copy):', dataURL);
                showNotification('Image data shown. Right-click to copy.', 'error');
            }
        }, 'image/png');
        
    } catch (error) {
        console.error('QR copy error:', error);
        showNotification('Failed to copy QR code. Try downloading instead.', 'error');
    }
}

function downloadQRCode() {
    if (!window.currentQRCodeUrl || !window.currentQRCodeData) {
        showNotification('No QR code available', 'error');
        return;
    }
    
    try {
        // Try to get image from modal or result
        let img = document.getElementById('modalQRImage');
        if (!img) {
            img = document.getElementById('qrCodeImage');
        }
        
        if (img && img.complete) {
            // Use canvas to ensure proper download
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            canvas.width = img.naturalWidth || img.width;
            canvas.height = img.naturalHeight || img.height;
            ctx.drawImage(img, 0, 0);
            
            // Convert to blob and download
            canvas.toBlob((blob) => {
                if (!blob) {
                    // Fallback to direct URL download
                    downloadFromURL();
                    return;
                }
                
                const url = URL.createObjectURL(blob);
                const link = document.createElement('a');
                link.href = url;
                link.download = `qrcode-${window.currentQRCodeData}.png`;
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);
                URL.revokeObjectURL(url);
                
                showNotification('QR Code downloaded!', 'success');
            }, 'image/png');
        } else {
            // Fallback: download directly from URL
            downloadFromURL();
        }
    } catch (error) {
        console.error('Download error:', error);
        downloadFromURL();
    }
}

function downloadFromURL() {
    const link = document.createElement('a');
    link.href = window.currentQRCodeUrl;
    link.download = `qrcode-${window.currentQRCodeData}.png`;
    link.target = '_blank';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    showNotification('QR Code download started!', 'success');
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showNotification('Copied to clipboard!', 'success');
    }).catch(() => {
        showNotification('Failed to copy', 'error');
    });
}

function showNotification(message, type) {
    // Simple notification (can be enhanced with a toast library)
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 20px;
        background: ${type === 'success' ? '#10b981' : '#ef4444'};
        color: white;
        border-radius: 8px;
        box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        z-index: 10000;
        animation: slideIn 0.3s ease;
    `;
    document.body.appendChild(notification);
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('authModal');
    if (event.target === modal) {
        closeAuthModal();
    }
}

