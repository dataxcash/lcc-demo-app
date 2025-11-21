const RuntimePage = {
    instanceId: null,
    productId: null,
    simulationRunning: false,
    statusUpdateInterval: null,
    eventUpdateInterval: null,
    chart: null,
    
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container page-enter-active">
                <div class="text-center mb-6">
                    <h1 style="font-size: var(--text-3xl); margin-bottom: var(--space-4); color: var(--text-secondary);">
                        📊 Live Simulation Dashboard
                    </h1>
                    <p style="font-size: var(--text-lg); color: var(--text-muted); max-width: 800px; margin: 0 auto;">
                        Watch the SDK in action as it executes feature checks and manages resource limits in real-time.
                    </p>
                </div>

                <div class="card mb-4">
                    <h3 class="card-title">Simulation Control</h3>
                    <div style="display: flex; gap: var(--space-3); align-items: center; margin-bottom: var(--space-4);">
                        <div style="flex: 1; min-width: 200px;">
                            <label style="font-weight: var(--weight-medium);">Instance:</label>
                            <select id="instance-select" class="form-input" style="width: 100%; margin-top: var(--space-2);">
                                <option value="">-- Select registered instance --</option>
                            </select>
                        </div>
                        <div>
                            <label style="font-weight: var(--weight-medium);">Status:</label>
                            <div id="status-badge" style="margin-top: var(--space-2); padding: var(--space-2) var(--space-3); background: var(--bg-soft); border-radius: var(--radius); color: var(--text-muted);">
                                Idle
                            </div>
                        </div>
                    </div>

                    <div style="display: flex; gap: var(--space-3); flex-wrap: wrap;">
                        <button id="btn-start" class="btn btn-primary" disabled>Start Simulation</button>
                        <button id="btn-pause" class="btn btn-secondary" disabled>Pause</button>
                        <button id="btn-stop" class="btn btn-danger" disabled>Stop</button>
                    </div>
                </div>

                <div class="card mb-4">
                    <h3 class="card-title">Progress</h3>
                    <div style="margin-bottom: var(--space-4);">
                        <div style="display: flex; justify-content: space-between; margin-bottom: var(--space-2);">
                            <span id="progress-text" style="color: var(--text-primary);">Ready</span>
                            <span id="time-info" style="color: var(--text-muted);">--</span>
                        </div>
                        <div style="width: 100%; height: 24px; background: var(--bg-soft); border-radius: var(--radius); overflow: hidden;">
                            <div id="progress-bar" style="width: 0%; height: 100%; background: linear-gradient(90deg, #3b82f6, #8b5cf6); transition: width 0.3s ease;"></div>
                        </div>
                    </div>
                </div>

                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: var(--space-4); margin-bottom: var(--space-4);">
                    <div class="card">
                        <h3 class="card-title">Statistics</h3>
                        <div style="display: grid; gap: var(--space-2);">
                            <div style="display: flex; justify-content: space-between; color: var(--text-primary);">
                                <span>Success:</span>
                                <span id="stat-success" style="font-weight: var(--weight-medium);">0</span>
                            </div>
                            <div style="display: flex; justify-content: space-between; color: var(--text-primary);">
                                <span>Failed:</span>
                                <span id="stat-failed" style="font-weight: var(--weight-medium); color: var(--error);">0</span>
                            </div>
                            <div style="display: flex; justify-content: space-between; color: var(--text-primary);">
                                <span>Success Rate:</span>
                                <span id="stat-rate" style="font-weight: var(--weight-medium);">--%</span>
                            </div>
                        </div>
                    </div>

                    <div class="card">
                        <h3 class="card-title">Events</h3>
                        <div style="display: grid; gap: var(--space-2);">
                            <div style="display: flex; justify-content: space-between; color: var(--text-primary);">
                                <span>Total Events:</span>
                                <span id="event-count" style="font-weight: var(--weight-medium);">0</span>
                            </div>
                            <div style="display: flex; justify-content: space-between; color: var(--text-primary);">
                                <span>Latest:</span>
                                <span id="event-latest" style="font-weight: var(--weight-medium); font-family: var(--font-mono); font-size: var(--text-sm);">-</span>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="card">
                    <h3 class="card-title">Event Log (Latest 20)</h3>
                    <div id="event-log" style="background: var(--bg-soft); border: 1px solid var(--border); border-radius: var(--radius); padding: var(--space-3); max-height: 400px; overflow-y: auto; font-family: var(--font-mono); font-size: var(--text-sm); color: var(--text-primary);">
                        <div style="color: var(--text-muted);">Waiting for simulation to start...</div>
                    </div>
                </div>

                <div class="card">
                    <h3 class="card-title">Quick Export</h3>
                    <button id="btn-export" class="btn btn-secondary" disabled>Export Results as JSON</button>
                </div>
            </div>
        `;

        await this.init();
    },

    async init() {
        await this.loadInstances();
        this.attachEventListeners();
    },

    async loadInstances() {
        try {
            const products = await Utils.fetchAPI('/api/sim/registered');
            const select = document.getElementById('instance-select');
            select.innerHTML = '<option value="">-- Select registered instance --</option>';
            
            if (Array.isArray(products)) {
                products.forEach(productId => {
                    const option = document.createElement('option');
                    option.value = productId;
                    option.textContent = productId;
                    select.appendChild(option);
                });
            }
        } catch (error) {
            console.error('Failed to load instances:', error);
        }
    },

    attachEventListeners() {
        const instanceSelect = document.getElementById('instance-select');
        instanceSelect.addEventListener('change', () => this.handleInstanceSelect());

        document.getElementById('btn-start').addEventListener('click', () => this.startSimulation());
        document.getElementById('btn-pause').addEventListener('click', () => this.pauseSimulation());
        document.getElementById('btn-stop').addEventListener('click', () => this.stopSimulation());
        document.getElementById('btn-export').addEventListener('click', () => this.exportResults());
    },

    handleInstanceSelect() {
        this.instanceId = document.getElementById('instance-select').value;
        this.productId = this.instanceId;
        
        const btnStart = document.getElementById('btn-start');
        btnStart.disabled = !this.instanceId;
    },

    async startSimulation() {
        if (!this.instanceId) {
            Utils.showAlert('warning', 'Please select an instance first');
            return;
        }

        try {
            const response = await Utils.fetchAPI('/api/simulation/start', {
                method: 'POST',
                body: JSON.stringify({
                    instance_id: this.instanceId,
                    iterations: 50,
                    interval_ms: 200,
                    features_to_call: ['basic_reports', 'ml_analytics', 'pdf_export'],
                    call_pattern: {
                        'basic_reports': 1,
                        'ml_analytics': 1,
                        'pdf_export': 3
                    }
                })
            });

            if (response.success) {
                this.simulationRunning = true;
                this.updateControlButtons();
                this.startUpdates();
                Utils.showAlert('success', 'Simulation started!');
            } else {
                Utils.showAlert('error', response.error || 'Failed to start simulation');
            }
        } catch (error) {
            Utils.showAlert('error', `Start failed: ${error.message}`);
        }
    },

    async pauseSimulation() {
        if (!this.instanceId) return;
        
        try {
            await Utils.fetchAPI(`/api/simulation/pause?instance_id=${this.instanceId}`, {
                method: 'POST'
            });
            this.simulationRunning = false;
            this.updateControlButtons();
        } catch (error) {
            Utils.showAlert('error', `Pause failed: ${error.message}`);
        }
    },

    async stopSimulation() {
        if (!this.instanceId) return;
        
        try {
            await Utils.fetchAPI(`/api/simulation/stop?instance_id=${this.instanceId}`, {
                method: 'POST'
            });
            this.simulationRunning = false;
            this.stopUpdates();
            this.updateControlButtons();
        } catch (error) {
            Utils.showAlert('error', `Stop failed: ${error.message}`);
        }
    },

    async exportResults() {
        if (!this.instanceId) return;
        
        try {
            const response = await Utils.fetchAPI(`/api/simulation/export?instance_id=${this.instanceId}`, {
                method: 'POST'
            });
            
            const dataStr = JSON.stringify(response, null, 2);
            const dataBlob = new Blob([dataStr], { type: 'application/json' });
            const url = URL.createObjectURL(dataBlob);
            const link = document.createElement('a');
            link.href = url;
            link.download = `simulation-${this.instanceId}-${new Date().getTime()}.json`;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            URL.revokeObjectURL(url);
            
            Utils.showAlert('success', 'Results exported successfully!');
        } catch (error) {
            Utils.showAlert('error', `Export failed: ${error.message}`);
        }
    },

    startUpdates() {
        this.stopUpdates();
        this.statusUpdateInterval = setInterval(() => this.updateStatus(), 500);
        this.eventUpdateInterval = setInterval(() => this.updateEvents(), 1000);
    },

    stopUpdates() {
        if (this.statusUpdateInterval) clearInterval(this.statusUpdateInterval);
        if (this.eventUpdateInterval) clearInterval(this.eventUpdateInterval);
    },

    async updateStatus() {
        try {
            const response = await Utils.fetchAPI(`/api/simulation/status?instance_id=${this.instanceId}`);
            
            if (response.success) {
                const metrics = response.metrics;
                const status = response.status;
                
                // Update status badge
                const badge = document.getElementById('status-badge');
                let statusColor = 'var(--text-muted)';
                if (status === 'running') statusColor = 'var(--success)';
                else if (status === 'paused') statusColor = 'var(--warning)';
                else if (status === 'stopped' || status === 'completed') statusColor = 'var(--error)';
                
                badge.textContent = status.toUpperCase();
                badge.style.color = statusColor;
                
                // Update progress
                const percent = metrics.total_iterations > 0 ? (metrics.completed_iterations / metrics.total_iterations) * 100 : 0;
                document.getElementById('progress-bar').style.width = percent + '%';
                
                const remaining = metrics.estimated_remaining_seconds > 0 ? metrics.estimated_remaining_seconds.toFixed(1) : '0';
                document.getElementById('progress-text').textContent = `${metrics.completed_iterations}/${metrics.total_iterations} iterations`;
                document.getElementById('time-info').textContent = `Elapsed: ${metrics.elapsed_seconds.toFixed(1)}s | Remaining: ~${remaining}s`;
                
                // Update statistics
                document.getElementById('stat-success').textContent = metrics.success_count;
                document.getElementById('stat-failed').textContent = metrics.failure_count;
                
                const total = metrics.success_count + metrics.failure_count;
                const rate = total > 0 ? ((metrics.success_count / total) * 100).toFixed(1) : '0';
                document.getElementById('stat-rate').textContent = rate + '%';
            }
        } catch (error) {
            console.error('Status update failed:', error);
        }
    },

    async updateEvents() {
        try {
            const response = await Utils.fetchAPI(`/api/simulation/events?instance_id=${this.instanceId}&limit=100`);
            
            if (response.success) {
                const events = response.events || [];
                document.getElementById('event-count').textContent = events.length;
                
                if (events.length > 0) {
                    const latest = events[events.length - 1];
                    const latestText = `${latest.type.replace('_', ' ')} - ${latest.feature_id || ''}`;
                    document.getElementById('event-latest').textContent = latestText.substring(0, 30);
                    
                    // Update event log
                    const log = document.getElementById('event-log');
                    const recentEvents = events.slice(-20).reverse();
                    
                    log.innerHTML = recentEvents.map(e => {
                        const color = e.allowed ? 'var(--success)' : 'var(--error)';
                        const icon = e.allowed ? '✓' : '✗';
                        return `<div style="color: ${color}; margin-bottom: var(--space-1);"><strong>${icon}</strong> [${e.iteration}] ${e.feature_id}: ${e.reason}</div>`;
                    }).join('');
                }
            }
        } catch (error) {
            console.error('Events update failed:', error);
        }
    },

    updateControlButtons() {
        const btnStart = document.getElementById('btn-start');
        const btnPause = document.getElementById('btn-pause');
        const btnStop = document.getElementById('btn-stop');
        const btnExport = document.getElementById('btn-export');
        
        btnStart.disabled = !this.instanceId || this.simulationRunning;
        btnPause.disabled = !this.simulationRunning;
        btnStop.disabled = !this.simulationRunning && !this.instanceId;
        btnExport.disabled = !this.instanceId;
    }
};
