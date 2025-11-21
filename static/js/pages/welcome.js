const WelcomePage = {
    savedConfig: null,
    currentInput: null,
    
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container page-enter-active">
                <div class="text-center mb-6">
                    <h1 style="font-size: var(--text-3xl); margin-bottom: var(--space-4); color: var(--text-secondary);">
                        Welcome to the LCC SDK Tutorial Platform!
                    </h1>
                    <p style="font-size: var(--text-lg); color: var(--text-muted); max-width: 800px; margin: 0 auto;">
                        This interactive demo will guide you through understanding license tiers, control types, 
                        SDK integration examples, and live runtime behavior observation.
                    </p>
                </div>

                <div class="card" style="max-width: 800px; margin: 0 auto;">
                    <h2 class="card-title">🔧 LCC Server Configuration</h2>
                    
                    <div id="config-alert" class="hidden"></div>
                    
                    <div class="form-group">
                        <label class="form-label">LCC URL</label>
                        <input 
                            type="text" 
                            id="lcc-url" 
                            class="form-input" 
                            placeholder="http://localhost:7086"
                            value=""
                        >
                    </div>

                    <div id="connection-status" class="hidden mb-4">
                        <div style="display: flex; align-items: center; gap: var(--space-3);">
                            <span style="font-weight: var(--weight-medium);">Status:</span>
                            <span id="status-badge"></span>
                        </div>
                    </div>

                    <div style="display: flex; gap: var(--space-3); margin-top: var(--space-6);">
                        <button id="btn-test" class="btn btn-secondary" style="flex: 1;">
                            Test Connection
                        </button>
                        <button id="btn-update" class="btn btn-primary" style="flex: 1;" disabled>
                            UPDATE
                        </button>
                        <button id="btn-continue" class="btn btn-primary" style="flex: 1;" disabled>
                            Save & Continue →
                        </button>
                    </div>
                </div>

                <div style="max-width: 800px; margin: var(--space-8) auto 0;">
                    <div class="card">
                        <p style="color: var(--text-muted); margin-bottom: var(--space-4);">
                            💡 First time here? The default URL works if LCC is running locally. 
                            Otherwise, update to your server.
                        </p>
                        
                        <h3 style="font-size: var(--text-lg); font-weight: var(--weight-semibold); 
                                   color: var(--text-secondary); margin-bottom: var(--space-3);">
                            📚 What You'll Learn:
                        </h3>
                        <ul style="list-style: none; padding: 0; color: var(--text-primary);">
                            <li style="margin-bottom: var(--space-2);">
                                • How licenses control feature availability (Tiers)
                            </li>
                            <li style="margin-bottom: var(--space-2);">
                                • How limits control usage amounts (Quota/TPS/etc)
                            </li>
                            <li style="margin-bottom: var(--space-2);">
                                • How to integrate SDK into your application
                            </li>
                            <li style="margin-bottom: var(--space-2);">
                                • How to handle license checks and denials gracefully
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        `;

        await this.init();
    },

    async init() {
        await this.loadConfig();
        this.attachEventListeners();
    },

    async loadConfig() {
        try {
            const data = await Utils.fetchAPI('/api/config');
            this.savedConfig = data.lcc_url || 'http://localhost:7086';
            this.currentInput = this.savedConfig;
            
            const urlInput = document.getElementById('lcc-url');
            urlInput.value = this.savedConfig;
            
            await this.testConnection();
        } catch (error) {
            console.error('Failed to load config:', error);
            this.currentInput = 'http://localhost:7086';
            document.getElementById('lcc-url').value = this.currentInput;
        }
    },

    attachEventListeners() {
        const urlInput = document.getElementById('lcc-url');
        const btnTest = document.getElementById('btn-test');
        const btnUpdate = document.getElementById('btn-update');
        const btnContinue = document.getElementById('btn-continue');

        urlInput.addEventListener('input', Utils.debounce(() => {
            this.currentInput = urlInput.value.trim();
            this.checkConfigChanged();
        }, 300));

        btnTest.addEventListener('click', () => this.testConnection());
        btnUpdate.addEventListener('click', () => this.updateConfig());
        btnContinue.addEventListener('click', () => this.continue());
    },

    checkConfigChanged() {
        const changed = this.currentInput !== this.savedConfig;
        const btnUpdate = document.getElementById('btn-update');
        const btnContinue = document.getElementById('btn-continue');
        const alertDiv = document.getElementById('config-alert');

        if (changed) {
            btnUpdate.disabled = false;
            btnContinue.disabled = true;
            alertDiv.className = 'alert alert-warning';
            alertDiv.textContent = '⚠️ Configuration changed - click UPDATE to apply';
            alertDiv.classList.remove('hidden');
        } else {
            btnUpdate.disabled = true;
            alertDiv.classList.add('hidden');
        }
    },

    async testConnection() {
        const btnTest = document.getElementById('btn-test');
        const statusDiv = document.getElementById('connection-status');
        const statusBadge = document.getElementById('status-badge');
        const btnContinue = document.getElementById('btn-continue');

        btnTest.disabled = true;
        btnTest.innerHTML = '<span class="spinner"></span> Testing...';

        try {
            const data = await Utils.fetchAPI('/api/config/validate');
            
            statusDiv.classList.remove('hidden');
            
            if (data.reachable) {
                statusBadge.className = 'badge badge-success';
                statusBadge.innerHTML = `✓ Connected | Version: ${data.version || 'Unknown'} | Products: ${data.products_count || 0}`;
                
                Utils.updateLCCStatus(true, { version: data.version });
                
                if (this.currentInput === this.savedConfig) {
                    btnContinue.disabled = false;
                }
            } else {
                statusBadge.className = 'badge badge-error';
                statusBadge.textContent = `✗ Failed: ${data.error || 'Connection failed'}`;
                Utils.updateLCCStatus(false);
                btnContinue.disabled = true;
            }
        } catch (error) {
            statusDiv.classList.remove('hidden');
            statusBadge.className = 'badge badge-error';
            statusBadge.textContent = `✗ Error: ${error.message}`;
            Utils.updateLCCStatus(false);
            btnContinue.disabled = true;
        } finally {
            btnTest.disabled = false;
            btnTest.textContent = 'Test Connection';
        }
    },

    async updateConfig() {
        const btnUpdate = document.getElementById('btn-update');
        const urlInput = document.getElementById('lcc-url');

        btnUpdate.disabled = true;
        btnUpdate.innerHTML = '<span class="spinner"></span> Updating...';

        try {
            await Utils.fetchAPI('/api/config', {
                method: 'POST',
                body: JSON.stringify({
                    lcc_url: this.currentInput
                })
            });

            this.savedConfig = this.currentInput;
            
            const alertDiv = document.getElementById('config-alert');
            alertDiv.className = 'alert alert-success';
            alertDiv.textContent = '✓ Configuration saved successfully';
            alertDiv.classList.remove('hidden');

            setTimeout(() => alertDiv.classList.add('hidden'), 3000);
            
            await this.testConnection();
        } catch (error) {
            Utils.showAlert('error', `Failed to save configuration: ${error.message}`);
        } finally {
            btnUpdate.disabled = true;
            btnUpdate.textContent = 'UPDATE';
        }
    },

    continue() {
        App.navigateTo('tiers');
    }
};
