const SetupMultiPage = {
    products: null,
    instances: [],  // List of registered instances
    pendingRegistrations: [],  // Registrations in progress
    
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container page-enter-active">
                <div class="text-center mb-6">
                    <h1 style="font-size: var(--text-3xl); margin-bottom: var(--space-4); color: var(--text-secondary);">
                        ⚙️ Lesson 3: SDK Instance Setup (Multi-Instance)
                    </h1>
                    <p style="font-size: var(--text-lg); color: var(--text-muted); max-width: 800px; margin: 0 auto;">
                        Register multiple SDK instances for different products and versions.
                        Manage them all from a single dashboard.
                    </p>
                </div>

                <!-- Instances Dashboard -->
                <div class="card mb-4">
                    <h3 class="card-title">📊 Registered Instances</h3>
                    <div id="instances-list" style="min-height: 200px;">
                        <div style="color: var(--text-muted); text-align: center; padding: var(--space-6);">
                            No instances registered yet. Create one below.
                        </div>
                    </div>
                </div>

                <!-- Batch Registration -->
                <div class="card mb-4">
                    <h3 class="card-title">➕ Register New Instance(s)</h3>
                    
                    <div style="display: grid; gap: var(--space-4); margin-bottom: var(--space-4);">
                        <div>
                            <label style="font-weight: var(--weight-medium); display: block; margin-bottom: var(--space-2);">
                                LCC Server URL:
                            </label>
                            <input type="text" id="lcc-url" class="form-input" 
                                   placeholder="http://localhost:7086"
                                   style="width: 100%; max-width: 500px;">
                        </div>

                        <div>
                            <label style="font-weight: var(--weight-medium); display: block; margin-bottom: var(--space-2);">
                                Select Products to Register:
                            </label>
                            <div id="products-checklist" style="display: grid; gap: var(--space-2); max-width: 600px;">
                                <div style="color: var(--text-muted);">Loading products...</div>
                            </div>
                        </div>

                        <div>
                            <label style="font-weight: var(--weight-medium); display: block; margin-bottom: var(--space-2);">
                                Version (applies to all selected products):
                            </label>
                            <input type="text" id="version-input" class="form-input" 
                                   value="1.0.0"
                                   style="width: 100%; max-width: 200px;">
                        </div>

                        <div style="display: flex; gap: var(--space-3); align-items: center;">
                            <button id="btn-register-batch" class="btn btn-primary">
                                Register Selected Products
                            </button>
                            <button id="btn-select-all" class="btn btn-secondary">Select All</button>
                            <button id="btn-clear-selection" class="btn btn-secondary">Clear Selection</button>
                        </div>
                    </div>

                    <div id="registration-progress" class="hidden" style="margin-top: var(--space-4);">
                        <div style="display: grid; gap: var(--space-2);" id="progress-items"></div>
                    </div>
                </div>

                <!-- Instance Management -->
                <div class="card mb-4">
                    <h3 class="card-title">🛠️ Instance Management</h3>
                    <p style="color: var(--text-muted); margin-bottom: var(--space-4);">
                        Test and manage your registered instances:
                    </p>

                    <div id="management-panel" style="display: grid; gap: var(--space-4);">
                        <div>
                            <label style="font-weight: var(--weight-medium); display: block; margin-bottom: var(--space-2);">
                                Select Instance:
                            </label>
                            <select id="instance-select" class="form-input" style="width: 100%; max-width: 400px;">
                                <option value="">-- No instances available --</option>
                            </select>
                        </div>

                        <div id="selected-instance-details" class="hidden" 
                             style="background: var(--bg-soft); padding: var(--space-4); border-radius: var(--radius);">
                            <div style="display: grid; grid-template-columns: 150px 1fr; gap: var(--space-2); color: var(--text-primary); margin-bottom: var(--space-4);">
                                <div style="font-weight: var(--weight-medium);">Instance ID:</div>
                                <div id="detail-instance-id" style="font-family: var(--font-mono); color: var(--accent);"></div>
                                
                                <div style="font-weight: var(--weight-medium);">Product ID:</div>
                                <div id="detail-product-id"></div>
                                
                                <div style="font-weight: var(--weight-medium);">Version:</div>
                                <div id="detail-version"></div>
                                
                                <div style="font-weight: var(--weight-medium);">Status:</div>
                                <div id="detail-status"></div>
                                
                                <div style="font-weight: var(--weight-medium);">Registered At:</div>
                                <div id="detail-registered-at"></div>
                            </div>

                            <div style="display: flex; gap: var(--space-3);">
                                <button id="btn-test-instance" class="btn btn-secondary">Test Connection</button>
                                <button id="btn-delete-instance" class="btn btn-danger">Delete Instance</button>
                            </div>

                            <div id="test-output" class="hidden" style="margin-top: var(--space-4);">
                                <h4 style="color: var(--text-secondary); margin-bottom: var(--space-2);">Test Output:</h4>
                                <pre id="test-output-json" style="background: var(--bg); padding: var(--space-3); border-radius: var(--radius); overflow: auto; max-height: 300px;"></pre>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Info -->
                <div class="card">
                    <h3 class="card-title">💡 Multi-Instance Benefits</h3>
                    <ul style="list-style: none; padding: 0; color: var(--text-primary);">
                        <li style="margin-bottom: var(--space-2);">• <strong>Test Multiple Products</strong> - Register different products in parallel</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>Version Testing</strong> - Try different versions of the same product</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>A/B Testing</strong> - Compare different instance configurations</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>Load Testing</strong> - Run stress tests with multiple instances</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>Instance Isolation</strong> - Each instance has its own ID and quota tracking</li>
                    </ul>
                </div>
            </div>
        `;

        await this.init();
    },

    async init() {
        await this.loadProducts();
        await this.loadInstances();
        await this.loadConfig();
        this.attachEventListeners();
    },

    async loadProducts() {
        try {
            this.products = await Utils.fetchAPI('/api/products');
            this.renderProductsChecklist();
        } catch (error) {
            console.error('Failed to load products:', error);
            Utils.showAlert('error', 'Failed to load products. Please configure LCC URL first.');
        }
    },

    async loadInstances() {
        try {
            const response = await Utils.fetchAPI('/api/instances');
            this.instances = response.instances || [];
            this.renderInstancesList();
            this.renderInstanceSelect();
        } catch (error) {
            console.error('Failed to load instances:', error);
        }
    },

    async loadConfig() {
        try {
            const config = await Utils.fetchAPI('/api/config');
            if (config.lcc_url) {
                document.getElementById('lcc-url').value = config.lcc_url;
            }
        } catch (error) {
            console.error('Failed to load config:', error);
        }
    },

    renderProductsChecklist() {
        const container = document.getElementById('products-checklist');
        if (!this.products || this.products.length === 0) {
            container.innerHTML = '<div style="color: var(--text-muted);">No products available</div>';
            return;
        }

        container.innerHTML = this.products.map(product => `
            <label style="display: flex; align-items: center; gap: var(--space-2); cursor: pointer; padding: var(--space-2); border-radius: var(--radius); transition: background 0.2s;">
                <input type="checkbox" class="product-checkbox" value="${product.id}" 
                       style="cursor: pointer; width: 18px; height: 18px;">
                <span style="font-weight: var(--weight-medium);">${product.name || product.id}</span>
            </label>
        `).join('');
    },

    renderInstancesList() {
        const container = document.getElementById('instances-list');
        if (this.instances.length === 0) {
            container.innerHTML = `
                <div style="color: var(--text-muted); text-align: center; padding: var(--space-6);">
                    No instances registered yet. Create one below.
                </div>
            `;
            return;
        }

        container.innerHTML = `
            <div style="display: grid; gap: var(--space-3);">
                ${this.instances.map((inst, idx) => `
                    <div style="background: var(--bg-soft); padding: var(--space-3); border-radius: var(--radius); border-left: 3px solid var(--accent);">
                        <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: var(--space-2); color: var(--text-primary);">
                            <div>
                                <div style="font-size: var(--text-sm); color: var(--text-muted);">Instance ID</div>
                                <div style="font-family: var(--font-mono); font-size: var(--text-sm); color: var(--accent);">${inst.instance_id}</div>
                            </div>
                            <div>
                                <div style="font-size: var(--text-sm); color: var(--text-muted);">Product</div>
                                <div style="font-weight: var(--weight-medium);">${inst.product_id}</div>
                            </div>
                            <div>
                                <div style="font-size: var(--text-sm); color: var(--text-muted);">Version</div>
                                <div>${inst.version}</div>
                            </div>
                            <div>
                                <div style="font-size: var(--text-sm); color: var(--text-muted);">Status</div>
                                <div style="color: var(--success);">✓ Active</div>
                            </div>
                        </div>
                    </div>
                `).join('')}
            </div>
        `;
    },

    renderInstanceSelect() {
        const select = document.getElementById('instance-select');
        if (this.instances.length === 0) {
            select.innerHTML = '<option value="">-- No instances available --</option>';
            return;
        }

        select.innerHTML = `
            <option value="">-- Select an instance --</option>
            ${this.instances.map((inst, idx) => `
                <option value="${idx}">${inst.product_id} v${inst.version} (${inst.instance_id.substring(0, 8)}...)</option>
            `).join('')}
        `;
    },

    attachEventListeners() {
        document.getElementById('btn-register-batch').addEventListener('click', () => this.registerBatch());
        document.getElementById('btn-select-all').addEventListener('click', () => this.selectAllProducts());
        document.getElementById('btn-clear-selection').addEventListener('click', () => this.clearProductSelection());
        
        document.getElementById('instance-select').addEventListener('change', (e) => this.handleInstanceSelect(e.target.value));
        document.getElementById('btn-test-instance').addEventListener('click', () => this.testInstance());
        document.getElementById('btn-delete-instance').addEventListener('click', () => this.deleteInstance());
    },

    getSelectedProducts() {
        return Array.from(document.querySelectorAll('.product-checkbox:checked'))
            .map(cb => cb.value);
    },

    selectAllProducts() {
        document.querySelectorAll('.product-checkbox').forEach(cb => cb.checked = true);
    },

    clearProductSelection() {
        document.querySelectorAll('.product-checkbox').forEach(cb => cb.checked = false);
    },

    async registerBatch() {
        const selectedProducts = this.getSelectedProducts();
        const lccUrl = document.getElementById('lcc-url').value;
        const version = document.getElementById('version-input').value || '1.0.0';

        if (selectedProducts.length === 0) {
            Utils.showAlert('warning', 'Please select at least one product');
            return;
        }

        if (!lccUrl) {
            Utils.showAlert('warning', 'Please enter LCC Server URL');
            return;
        }

        const progressDiv = document.getElementById('registration-progress');
        const progressItems = document.getElementById('progress-items');
        progressDiv.classList.remove('hidden');
        progressItems.innerHTML = selectedProducts.map(pid => `
            <div class="progress-item" data-product="${pid}">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: var(--space-2);">
                    <span>${pid}</span>
                    <span style="color: var(--text-muted);">Pending...</span>
                </div>
                <div style="background: var(--bg-soft); height: 6px; border-radius: 3px; overflow: hidden;">
                    <div class="progress-bar" style="width: 0%; height: 100%; background: var(--accent); transition: width 0.3s;"></div>
                </div>
            </div>
        `).join('');

        for (const productId of selectedProducts) {
            await this.registerProduct(productId, version, lccUrl);
        }

        await this.loadInstances();
        Utils.showAlert('success', `Registered ${selectedProducts.length} instance(s)`);
    },

    async registerProduct(productId, version, lccUrl) {
        const progressItem = document.querySelector(`[data-product="${productId}"]`);
        
        try {
            const response = await Utils.fetchAPI('/api/instance/register', {
                method: 'POST',
                body: JSON.stringify({
                    product_id: productId,
                    version: version,
                    lcc_url: lccUrl
                })
            });

            if (response.success) {
                progressItem.querySelector('span:last-child').textContent = '✓ Success';
                progressItem.querySelector('.progress-bar').style.width = '100%';
                this.instances.push(response);
            } else {
                progressItem.querySelector('span:last-child').textContent = '✗ Failed';
            }
        } catch (error) {
            progressItem.querySelector('span:last-child').textContent = `✗ Error: ${error.message}`;
        }
    },

    handleInstanceSelect(idx) {
        if (idx === '') {
            document.getElementById('selected-instance-details').classList.add('hidden');
            return;
        }

        const inst = this.instances[parseInt(idx)];
        if (!inst) return;

        document.getElementById('detail-instance-id').textContent = inst.instance_id;
        document.getElementById('detail-product-id').textContent = inst.product_id;
        document.getElementById('detail-version').textContent = inst.version;
        document.getElementById('detail-status').innerHTML = '<span style="color: var(--success);">✓ Active</span>';
        document.getElementById('detail-registered-at').textContent = new Date(inst.registered_at).toLocaleString();
        
        document.getElementById('selected-instance-details').classList.remove('hidden');
        document.getElementById('test-output').classList.add('hidden');
    },

    async testInstance() {
        const select = document.getElementById('instance-select');
        const idx = parseInt(select.value);
        if (idx < 0 || idx >= this.instances.length) return;

        const inst = this.instances[idx];
        
        try {
            const response = await Utils.fetchAPI('/api/instance/test', {
                method: 'POST',
                body: JSON.stringify({
                    product_id: inst.product_id,
                    instance_id: inst.instance_id
                })
            });

            const testOutput = document.getElementById('test-output');
            const outputJson = document.getElementById('test-output-json');
            outputJson.textContent = JSON.stringify(response, null, 2);
            testOutput.classList.remove('hidden');
        } catch (error) {
            Utils.showAlert('error', `Test failed: ${error.message}`);
        }
    },

    async deleteInstance() {
        const select = document.getElementById('instance-select');
        const idx = parseInt(select.value);
        if (idx < 0 || idx >= this.instances.length) return;

        if (!confirm('Are you sure you want to delete this instance?')) return;

        const inst = this.instances[idx];
        
        try {
            const response = await Utils.fetchAPI('/api/instance/clear', {
                method: 'POST',
                body: JSON.stringify({
                    product_id: inst.product_id,
                    instance_id: inst.instance_id
                })
            });

            if (response.success) {
                this.instances.splice(idx, 1);
                await this.loadInstances();
                document.getElementById('selected-instance-details').classList.add('hidden');
                Utils.showAlert('success', 'Instance deleted');
            }
        } catch (error) {
            Utils.showAlert('error', `Delete failed: ${error.message}`);
        }
    }
};
