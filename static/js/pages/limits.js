const LimitsPage = {
    limitTypes: null,
    currentType: 'quota',
    
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container page-enter-active">
                <div class="text-center mb-6">
                    <h1 style="font-size: var(--text-3xl); margin-bottom: var(--space-4); color: var(--text-secondary);">
                        📚 Lesson 2: Understanding License Limits
                    </h1>
                    <p style="font-size: var(--text-lg); color: var(--text-muted); max-width: 800px; margin: 0 auto;">
                        Beyond ON/OFF feature gating, limits control HOW MUCH customers can consume.
                        Four types serve different use cases.
                    </p>
                </div>

                <div class="card mb-4">
                    <div style="display: flex; gap: var(--space-2); margin-bottom: var(--space-4); flex-wrap: wrap;">
                        <button class="limit-type-tab active" data-type="quota">1️⃣ Quota</button>
                        <button class="limit-type-tab" data-type="tps">2️⃣ TPS</button>
                        <button class="limit-type-tab" data-type="capacity">3️⃣ Capacity</button>
                        <button class="limit-type-tab" data-type="concurrency">4️⃣ Concurrency</button>
                    </div>

                    <div id="limit-type-content" style="min-height: 600px; transition: opacity 0.2s ease-in-out;"></div>
                </div>

                <div class="card">
                    <h3 class="card-title">Control Type Comparison</h3>
                    <div id="comparison-table"></div>
                </div>
            </div>
        `;

        await this.init();
    },

    async init() {
        await this.loadLimitTypes();
        this.attachEventListeners();
        this.updateDisplay();
        this.renderComparisonTable();
    },

    async loadLimitTypes() {
        try {
            this.limitTypes = await Utils.fetchAPI('/api/limits/types');
        } catch (error) {
            console.error('Failed to load limit types:', error);
            Utils.showAlert('error', 'Failed to load limit types');
        }
    },

    attachEventListeners() {
        const tabs = document.querySelectorAll('.limit-type-tab');
        tabs.forEach(tab => {
            tab.addEventListener('click', () => {
                this.currentType = tab.dataset.type;
                this.updateDisplay();
            });
        });
    },

    async updateDisplay() {
        this.updateTabs();
        await this.loadAndRenderContent();
    },

    updateTabs() {
        document.querySelectorAll('.limit-type-tab').forEach(tab => {
            if (tab.dataset.type === this.currentType) {
                tab.classList.add('active');
            } else {
                tab.classList.remove('active');
            }
        });
    },

    async loadAndRenderContent() {
        const content = document.getElementById('limit-type-content');
        content.innerHTML = '<div style="text-align: center; padding: var(--space-4); color: var(--text-muted);">Loading...</div>';

        try {
            const typeInfo = this.limitTypes.find(t => t.type === this.currentType);
            const example = await Utils.fetchAPI(`/api/limits/${this.currentType}/example`);
            
            this.renderTypeContent(typeInfo, example);
        } catch (error) {
            console.error('Failed to load example:', error);
            content.innerHTML = '<div style="text-align: center; padding: var(--space-4); color: var(--error);">Failed to load content</div>';
        }
    },

    renderTypeContent(typeInfo, example) {
        const content = document.getElementById('limit-type-content');
        
        let html = `
            <div style="margin-bottom: var(--space-6);">
                <h2 style="font-size: var(--text-2xl); color: var(--text-secondary); margin-bottom: var(--space-3);">
                    ${typeInfo.title}
                </h2>

                <div style="background: var(--bg-soft); padding: var(--space-4); border-radius: var(--radius); margin-bottom: var(--space-4);">
                    <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-2);">
                        📖 What Is It?
                    </h4>
                    <p style="color: var(--text-primary); margin-bottom: var(--space-3);">
                        ${typeInfo.description}
                    </p>

                    <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-2);">
                        🎯 Use Cases:
                    </h4>
                    <ul style="color: var(--text-primary); margin-bottom: var(--space-3); padding-left: var(--space-5);">
                        ${typeInfo.use_cases.map(uc => `<li>${uc}</li>`).join('')}
                    </ul>

                    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: var(--space-4); margin-bottom: var(--space-3);">
                        <div>
                            <h4 style="font-size: var(--text-base); color: var(--text-secondary); margin-bottom: var(--space-1);">
                                ⏰ Time Dimension:
                            </h4>
                            <p style="color: var(--text-primary);">${typeInfo.time_dimension}</p>
                        </div>
                        <div>
                            <h4 style="font-size: var(--text-base); color: var(--text-secondary); margin-bottom: var(--space-1);">
                                🔢 Who Tracks?
                            </h4>
                            <p style="color: var(--text-primary);">${typeInfo.who_tracks}</p>
                        </div>
                    </div>

                    <div>
                        <h4 style="font-size: var(--text-base); color: var(--text-secondary); margin-bottom: var(--space-1);">
                            💻 SDK API:
                        </h4>
                        <code style="background: #0d1117; padding: var(--space-1) var(--space-2); border-radius: 4px; color: var(--accent);">
                            ${typeInfo.sdk_api}
                        </code>
                    </div>
                </div>

                <div style="margin-bottom: var(--space-4);">
                    <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-3);">
                        📄 License File Configuration:
                    </h4>
                    <pre class="code-block">${this.escapeHtml(example.license_config)}</pre>
                </div>

                <div style="margin-bottom: var(--space-4);">
                    <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-3);">
                        💻 SDK Integration Code:
                    </h4>
                    <pre class="code-block">${this.escapeHtml(example.code_example)}</pre>
                </div>

                <div style="margin-bottom: var(--space-4);">
                    <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-3);">
                        🔄 Runtime Behavior:
                    </h4>
                    <table style="width: 100%; border-collapse: collapse;">
                        <thead>
                            <tr style="border-bottom: 2px solid var(--border);">
                                <th style="padding: var(--space-2); text-align: left; color: var(--text-secondary);">Call #</th>
                                <th style="padding: var(--space-2); text-align: center; color: var(--text-secondary);">Allowed</th>
                                <th style="padding: var(--space-2); text-align: center; color: var(--text-secondary);">Remaining</th>
                                <th style="padding: var(--space-2); text-align: left; color: var(--text-secondary);">Reason</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${example.behavior_table.map(row => `
                                <tr style="border-bottom: 1px solid var(--border);">
                                    <td style="padding: var(--space-2); color: var(--text-primary);">${row.call}</td>
                                    <td style="padding: var(--space-2); text-align: center; color: ${row.allowed.includes('✓') ? 'var(--success)' : 'var(--error)'};">${row.allowed}</td>
                                    <td style="padding: var(--space-2); text-align: center; color: var(--text-muted);">${row.remaining}</td>
                                    <td style="padding: var(--space-2); color: var(--text-muted);">${row.reason}</td>
                                </tr>
                            `).join('')}
                        </tbody>
                    </table>
                </div>

                <div style="margin-bottom: var(--space-4);">
                    <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-3);">
                        ✨ Key Points:
                    </h4>
                    <ul style="color: var(--text-primary); padding-left: var(--space-5);">
                        ${example.key_points.map(kp => `<li style="margin-bottom: var(--space-2);">• ${kp}</li>`).join('')}
                    </ul>
                </div>

                <div class="card" style="background: var(--bg-soft); padding: var(--space-4);">
                    <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-3);">
                        🎮 Interactive Simulation
                    </h4>
                    <p style="color: var(--text-muted); margin-bottom: var(--space-3);">
                        Run a mini simulation to see how this limit type behaves:
                    </p>
                    
                    <div style="display: flex; gap: var(--space-3); margin-bottom: var(--space-3); align-items: center; flex-wrap: wrap;">
                        <label style="color: var(--text-primary);">Iterations:</label>
                        <input type="number" id="sim-iterations" class="form-input" value="10" min="1" max="50" style="width: 100px;">
                        <button id="btn-run-simulation" class="btn btn-primary">▶ Run Simulation</button>
                        <button id="btn-reset-simulation" class="btn btn-secondary">🔄 Reset</button>
                    </div>

                    <div id="simulation-results" class="hidden">
                        <h5 style="color: var(--text-secondary); margin-bottom: var(--space-2);">Results:</h5>
                        <div id="simulation-progress" style="margin-bottom: var(--space-3);"></div>
                        <div id="simulation-summary" style="padding: var(--space-3); background: #0d1117; border-radius: var(--radius); margin-bottom: var(--space-2); color: var(--text-primary);"></div>
                        <div style="max-height: 300px; overflow-y: auto; background: #0d1117; border-radius: var(--radius); padding: var(--space-3);">
                            <table id="simulation-table" style="width: 100%; border-collapse: collapse; font-size: var(--text-sm);">
                                <thead>
                                    <tr style="border-bottom: 1px solid var(--border);">
                                        <th style="padding: var(--space-1); text-align: left; color: var(--text-muted);">#</th>
                                        <th style="padding: var(--space-1); text-align: center; color: var(--text-muted);">Status</th>
                                        <th style="padding: var(--space-1); text-align: center; color: var(--text-muted);">Remaining</th>
                                        <th style="padding: var(--space-1); text-align: left; color: var(--text-muted);">Reason</th>
                                        <th style="padding: var(--space-1); text-align: left; color: var(--text-muted);">Details</th>
                                    </tr>
                                </thead>
                                <tbody id="simulation-tbody"></tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        `;

        content.innerHTML = html;
        this.attachSimulationListeners();
    },

    attachSimulationListeners() {
        const btnRun = document.getElementById('btn-run-simulation');
        const btnReset = document.getElementById('btn-reset-simulation');

        if (btnRun) {
            btnRun.addEventListener('click', () => this.runSimulation());
        }

        if (btnReset) {
            btnReset.addEventListener('click', () => this.resetSimulation());
        }
    },

    async runSimulation() {
        const iterations = parseInt(document.getElementById('sim-iterations').value) || 10;
        const resultsDiv = document.getElementById('simulation-results');
        const tbody = document.getElementById('simulation-tbody');
        const summary = document.getElementById('simulation-summary');

        resultsDiv.classList.remove('hidden');
        tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; padding: var(--space-3); color: var(--text-muted);">Running simulation...</td></tr>';

        try {
            const params = this.getSimulationParams();
            const response = await Utils.fetchAPI(`/api/limits/${this.currentType}/simulate`, {
                method: 'POST',
                body: JSON.stringify({
                    feature_id: this.currentType + '_feature',
                    iterations: iterations,
                    params: params
                })
            });

            if (response.success) {
                this.renderSimulationResults(response);
            } else {
                tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; padding: var(--space-3); color: var(--error);">Simulation failed</td></tr>';
            }
        } catch (error) {
            console.error('Simulation error:', error);
            tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; padding: var(--space-3); color: var(--error);">Error running simulation</td></tr>';
        }
    },

    getSimulationParams() {
        switch (this.currentType) {
            case 'quota':
                return { max: 10000, amount: 1 };
            case 'tps':
                return { max_tps: 10.0 };
            case 'capacity':
                return { max_capacity: 50 };
            case 'concurrency':
                return { max_concurrency: 10 };
            default:
                return {};
        }
    },

    renderSimulationResults(response) {
        const tbody = document.getElementById('simulation-tbody');
        const summary = document.getElementById('simulation-summary');

        summary.textContent = response.summary;

        let html = '';
        response.results.forEach(result => {
            const statusColor = result.allowed ? 'var(--success)' : 'var(--error)';
            const statusIcon = result.allowed ? '✓' : '✗';
            
            html += `
                <tr style="border-bottom: 1px solid var(--border-soft);">
                    <td style="padding: var(--space-1); color: var(--text-muted);">${result.iteration}</td>
                    <td style="padding: var(--space-1); text-align: center; color: ${statusColor};">${statusIcon}</td>
                    <td style="padding: var(--space-1); text-align: center; color: var(--text-muted); font-family: var(--font-mono); font-size: 12px;">${result.remaining}</td>
                    <td style="padding: var(--space-1); color: var(--text-muted);">${result.reason}</td>
                    <td style="padding: var(--space-1); color: var(--text-muted); font-size: 12px;">${result.details || '-'}</td>
                </tr>
            `;
        });

        tbody.innerHTML = html;
    },

    resetSimulation() {
        const resultsDiv = document.getElementById('simulation-results');
        resultsDiv.classList.add('hidden');
        document.getElementById('sim-iterations').value = '10';
    },

    renderComparisonTable() {
        const tableDiv = document.getElementById('comparison-table');
        
        if (!this.limitTypes) {
            tableDiv.innerHTML = '<p style="color: var(--text-muted);">Loading...</p>';
            return;
        }

        let html = `
            <table style="width: 100%; border-collapse: collapse; margin-top: var(--space-3);">
                <thead>
                    <tr style="border-bottom: 2px solid var(--border);">
                        <th style="padding: var(--space-3); text-align: left; color: var(--text-secondary);">Control</th>
                        <th style="padding: var(--space-3); text-align: left; color: var(--text-secondary);">License Config</th>
                        <th style="padding: var(--space-3); text-align: left; color: var(--text-secondary);">SDK API</th>
                        <th style="padding: var(--space-3); text-align: left; color: var(--text-secondary);">Developer</th>
                    </tr>
                </thead>
                <tbody>
        `;

        const comparisonData = [
            {
                control: 'Quota',
                config: 'quota: {max, window}',
                api: 'Consume(id, amount)',
                developer: 'Just call Consume()'
            },
            {
                control: 'TPS',
                config: 'max_tps: 10.0',
                api: 'CheckTPS(id, current)',
                developer: 'Measure TPS'
            },
            {
                control: 'Capacity',
                config: 'max_capacity: 50',
                api: 'CheckCapacity(id, used)',
                developer: 'Count resources'
            },
            {
                control: 'Concurrency',
                config: 'max_concurrency: 10',
                api: 'AcquireSlot(id)',
                developer: 'Call + release()'
            }
        ];

        comparisonData.forEach(row => {
            html += `
                <tr style="border-bottom: 1px solid var(--border);">
                    <td style="padding: var(--space-3); color: var(--text-primary); font-weight: var(--weight-medium);">${row.control}</td>
                    <td style="padding: var(--space-3); color: var(--text-muted); font-family: var(--font-mono); font-size: var(--text-sm);">${row.config}</td>
                    <td style="padding: var(--space-3); color: var(--text-muted); font-family: var(--font-mono); font-size: var(--text-sm);">${row.api}</td>
                    <td style="padding: var(--space-3); color: var(--text-muted);">${row.developer}</td>
                </tr>
            `;
        });

        html += '</tbody></table>';
        tableDiv.innerHTML = html;
    },

    escapeHtml(text) {
        const map = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#039;'
        };
        return text.replace(/[&<>"']/g, m => map[m]);
    }
};
