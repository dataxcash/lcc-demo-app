const SetupPage = {
    products: null,
    selectedProductId: null,
    currentStep: 1,
    instanceData: null,
    
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container page-enter-active">
                <div class="text-center mb-6">
                    <h1 style="font-size: var(--text-3xl); margin-bottom: var(--space-4); color: var(--text-secondary);">
                        ⚙️ Lesson 3: SDK Instance Setup
                    </h1>
                    <p style="font-size: var(--text-lg); color: var(--text-muted); max-width: 800px; margin: 0 auto;">
                        Learn how to initialize and register an LCC SDK client instance.
                        This is the foundation for all SDK operations.
                    </p>
                </div>

                <div class="card mb-4">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: var(--space-4);">
                        <h3 class="card-title" style="margin: 0;">Step 1: Select Product</h3>
                        <button id="btn-view-instances" class="btn btn-secondary" style="font-size: var(--text-sm);">📊 View All Instances</button>
                    </div>
                    <p style="color: var(--text-muted); margin-bottom: var(--space-4);">
                        Choose which product you want to register an SDK instance for:
                    </p>
                    
                    <div style="display: flex; gap: var(--space-3); align-items: center; margin-bottom: var(--space-4);">
                        <label style="font-weight: var(--weight-medium); min-width: 120px;">Product:</label>
                        <select id="product-select" class="form-input" style="flex: 1; max-width: 400px;">
                            <option value="">-- Loading products... --</option>
                        </select>
                        <button id="btn-load-products" class="btn btn-secondary">Refresh Products</button>
                    </div>
                    
                    <!-- Multi-Instance Panel -->
                    <div id="instances-panel" class="hidden" style="margin-top: var(--space-4); padding: var(--space-4); background: var(--bg-soft); border-radius: var(--radius); border: 1px solid var(--border);">
                        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: var(--space-3);">
                            <h4 style="color: var(--text-secondary); margin: 0;">📋 Registered Instances</h4>
                            <button id="btn-close-instances" class="btn btn-secondary" style="font-size: var(--text-sm);">Close</button>
                        </div>
                        <div id="instances-list" style="max-height: 400px; overflow-y: auto;"></div>
                    </div>

                    <div id="product-info" class="hidden" style="background: var(--bg-soft); padding: var(--space-4); border-radius: var(--radius); margin-top: var(--space-4);">
                        <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-2);">Product Details:</h4>
                        <div id="product-details"></div>
                    </div>
                </div>

                <div class="card mb-4">
                    <h3 class="card-title">Step 2: Configure SDK Client</h3>
                    <p style="color: var(--text-muted); margin-bottom: var(--space-4);">
                        Set up the SDK configuration parameters:
                    </p>

                    <div style="display: grid; gap: var(--space-4);">
                        <div style="display: flex; gap: var(--space-3); align-items: center;">
                            <label style="font-weight: var(--weight-medium); min-width: 150px;">Product ID:</label>
                            <input type="text" id="config-product-id" class="form-input" style="flex: 1; max-width: 400px;" readonly placeholder="Select a product first">
                        </div>

                        <div style="display: flex; gap: var(--space-3); align-items: center;">
                            <label style="font-weight: var(--weight-medium); min-width: 150px;">Product Version:</label>
                            <input type="text" id="config-version" class="form-input" style="flex: 1; max-width: 400px;" value="1.0.0">
                        </div>

                        <div style="display: flex; gap: var(--space-3); align-items: center;">
                            <label style="font-weight: var(--weight-medium); min-width: 150px;">LCC Server URL:</label>
                            <input type="text" id="config-lcc-url" class="form-input" style="flex: 1; max-width: 400px;" placeholder="http://localhost:7086">
                        </div>
                    </div>

                    <div style="margin-top: var(--space-4); background: var(--bg-soft); padding: var(--space-4); border-radius: var(--radius);">
                        <h4 style="font-size: var(--text-base); color: var(--text-secondary); margin-bottom: var(--space-2);">
                            💡 What is SDK Configuration?
                        </h4>
                        <p style="color: var(--text-primary); margin: 0;">
                            The SDK needs to know which product it's managing and where the LCC server is located.
                            These parameters are passed to <code>client.NewClient()</code> during initialization.
                        </p>
                    </div>
                </div>

                <div class="card mb-4">
                    <h3 class="card-title">Step 3: Authentication Keys</h3>
                    <p style="color: var(--text-muted); margin-bottom: var(--space-4);">
                        The SDK uses RSA key pairs for secure authentication with the LCC server:
                    </p>

                    <div style="display: flex; gap: var(--space-3); margin-bottom: var(--space-4);">
                        <button id="btn-generate-keys" class="btn btn-secondary">Generate New Key Pair</button>
                        <button id="btn-use-saved-keys" class="btn btn-secondary">Use Saved Keys</button>
                    </div>

                    <div id="keys-status" style="color: var(--text-muted); padding: var(--space-2); background: var(--bg-soft); border-radius: var(--radius);">
                        Keys will be automatically generated or loaded from keystore during registration.
                    </div>

                    <div style="margin-top: var(--space-4); background: var(--bg-soft); padding: var(--space-4); border-radius: var(--radius);">
                        <h4 style="font-size: var(--text-base); color: var(--text-secondary); margin-bottom: var(--space-2);">
                            🔐 Authentication Flow:
                        </h4>
                        <ol style="color: var(--text-primary); padding-left: var(--space-5);">
                            <li>SDK generates RSA key pair (or loads from keystore)</li>
                            <li>Public key is sent to LCC server during registration</li>
                            <li>LCC server validates and creates instance ID</li>
                            <li>SDK signs all future requests with private key</li>
                        </ol>
                    </div>
                </div>

                <div class="card mb-4">
                    <h3 class="card-title">Step 4: Register Instance</h3>
                    <p style="color: var(--text-muted); margin-bottom: var(--space-4);">
                        Ready to register? This will call <code>client.Register()</code> on the LCC server:
                    </p>

                    <div style="display: flex; gap: var(--space-3); margin-bottom: var(--space-4);">
                        <button id="btn-register" class="btn btn-primary" disabled>Register Instance</button>
                        <button id="btn-clear-instance" class="btn btn-danger hidden">Clear Instance</button>
                    </div>

                    <div id="registration-status" class="hidden" style="margin-top: var(--space-4);"></div>

                    <div id="registration-code" style="margin-top: var(--space-4);">
                        <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-3);">
                            💻 What Happens During Registration:
                        </h4>
                        <pre class="code-block" style="margin-bottom: var(--space-4);">// SDK initialization and registration
cfg := &config.SDKConfig{
    LCCURL:         "http://localhost:7086",
    ProductID:      "data-insight-pro",
    ProductVersion: "1.0.0",
    Timeout:        10 * time.Second,
    CacheTTL:       5 * time.Second,
}

// Create client with key pair
keyPair, _ := auth.GenerateKeyPair()
client, err := client.NewClientWithKeyPair(cfg, keyPair)
if err != nil {
    return fmt.Errorf("client creation: %w", err)
}

// Register with LCC server
if err := client.Register(); err != nil {
    return fmt.Errorf("registration: %w", err)
}

// Now ready to use
instanceID := client.GetInstanceID()
log.Printf("Registered: %s", instanceID)</pre>
                    </div>
                </div>

                <div id="instance-status-card" class="card hidden">
                    <h3 class="card-title">Instance Status</h3>
                    
                    <div id="instance-info" style="background: var(--bg-soft); padding: var(--space-4); border-radius: var(--radius); margin-bottom: var(--space-4);">
                        <div style="display: grid; grid-template-columns: 150px 1fr; gap: var(--space-2); color: var(--text-primary);">
                            <div style="font-weight: var(--weight-medium);">Instance ID:</div>
                            <div id="instance-id-display" style="font-family: var(--font-mono); color: var(--accent);"></div>
                            
                            <div style="font-weight: var(--weight-medium);">Product ID:</div>
                            <div id="product-id-display"></div>
                            
                            <div style="font-weight: var(--weight-medium);">Status:</div>
                            <div id="status-display"></div>
                            
                            <div style="font-weight: var(--weight-medium);">Registered At:</div>
                            <div id="registered-at-display"></div>
                        </div>
                    </div>

                    <h4 style="font-size: var(--text-lg); color: var(--text-secondary); margin-bottom: var(--space-3);">
                        Test Instance Connection:
                    </h4>
                    
                    <div style="display: flex; gap: var(--space-3); align-items: center; margin-bottom: var(--space-4);">
                        <label style="font-weight: var(--weight-medium);">Feature to Test:</label>
                        <select id="test-feature-select" class="form-input" style="width: 250px;">
                            <option value="basic_reports">basic_reports</option>
                            <option value="ml_analytics">ml_analytics</option>
                            <option value="pdf_export">pdf_export</option>
                            <option value="excel_export">excel_export</option>
                        </select>
                        <button id="btn-test-feature" class="btn btn-primary">Test CheckFeature()</button>
                    </div>

                    <div id="test-result" class="hidden" style="margin-top: var(--space-4);">
                        <h4 style="font-size: var(--text-base); color: var(--text-secondary); margin-bottom: var(--space-2);">Test Result:</h4>
                        <pre id="test-result-json" class="code-block"></pre>
                    </div>
                </div>

                <div class="card">
                    <h3 class="card-title">📝 Key Concepts</h3>
                    <ul style="list-style: none; padding: 0; color: var(--text-primary);">
                        <li style="margin-bottom: var(--space-2);">• <strong>SDK Client</strong> = Your application's interface to LCC</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>Registration</strong> = Establishing trusted connection with LCC server</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>Instance ID</strong> = Unique identifier for this SDK instance</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>Key Pair</strong> = RSA public/private keys for secure authentication</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>Keystore</strong> = Local persistent storage for keys (reusable across restarts)</li>
                        <li style="margin-bottom: var(--space-2);">• <strong>Product ID</strong> = Which product this instance is licensed for</li>
                    </ul>
                </div>
            </div>
        `;

        await this.init();
    },

    async init() {
        await this.loadProducts();
        await this.loadConfig();
        this.attachEventListeners();
    },

    async loadProducts() {
        try {
            this.products = await Utils.fetchAPI('/api/products');
            this.renderProductSelect();
        } catch (error) {
            console.error('Failed to load products:', error);
            Utils.showAlert('error', 'Failed to load products. Please configure LCC URL first.');
        }
    },

    async loadConfig() {
        try {
            const config = await Utils.fetchAPI('/api/config');
            if (config.lcc_url) {
                document.getElementById('config-lcc-url').value = config.lcc_url;
            }
        } catch (error) {
            console.error('Failed to load config:', error);
        }
    },

    renderProductSelect() {
        const select = document.getElementById('product-select');
        select.innerHTML = '<option value="">-- Select a product --</option>';

        if (this.products && this.products.length > 0) {
            this.products.forEach(product => {
                const option = document.createElement('option');
                option.value = product.id;
                option.textContent = product.name || product.id;
                select.appendChild(option);
            });
        } else {
            select.innerHTML = '<option value="">No products available</option>';
        }
    },

    attachEventListeners() {
        const productSelect = document.getElementById('product-select');
        productSelect.addEventListener('change', () => this.handleProductSelect());

        const btnLoadProducts = document.getElementById('btn-load-products');
        btnLoadProducts.addEventListener('click', () => this.loadProducts());

        const btnRegister = document.getElementById('btn-register');
        btnRegister.addEventListener('click', () => this.registerInstance());

        const btnClear = document.getElementById('btn-clear-instance');
        btnClear.addEventListener('click', () => this.clearInstance());

        const btnTestFeature = document.getElementById('btn-test-feature');
        btnTestFeature.addEventListener('click', () => this.testFeature());

        const btnGenerateKeys = document.getElementById('btn-generate-keys');
        btnGenerateKeys.addEventListener('click', () => this.generateKeys());

        const btnUseSavedKeys = document.getElementById('btn-use-saved-keys');
        btnUseSavedKeys.addEventListener('click', () => this.useSavedKeys());

        // Multi-instance buttons
        const btnViewInstances = document.getElementById('btn-view-instances');
        btnViewInstances.addEventListener('click', () => this.viewAllInstances());

        const btnCloseInstances = document.getElementById('btn-close-instances');
        btnCloseInstances.addEventListener('click', () => this.closeInstancesPanel());
    },

    handleProductSelect() {
        const select = document.getElementById('product-select');
        const productId = select.value;

        if (!productId) {
            document.getElementById('product-info').classList.add('hidden');
            document.getElementById('config-product-id').value = '';
            document.getElementById('btn-register').disabled = true;
            return;
        }

        this.selectedProductId = productId;
        const product = this.products.find(p => p.id === productId);

        document.getElementById('config-product-id').value = productId;
        document.getElementById('btn-register').disabled = false;

        const productInfo = document.getElementById('product-info');
        productInfo.classList.remove('hidden');

        const details = document.getElementById('product-details');
        details.innerHTML = `
            <div style="display: grid; grid-template-columns: 120px 1fr; gap: var(--space-2); color: var(--text-primary);">
                <div style="font-weight: var(--weight-medium);">Product ID:</div>
                <div style="font-family: var(--font-mono);">${product.id}</div>
                
                <div style="font-weight: var(--weight-medium);">Name:</div>
                <div>${product.name || 'N/A'}</div>
                
                <div style="font-weight: var(--weight-medium);">Description:</div>
                <div>${product.description || 'N/A'}</div>
            </div>
        `;
    },

    async registerInstance() {
        const productId = document.getElementById('config-product-id').value;
        const version = document.getElementById('config-version').value;
        const lccUrl = document.getElementById('config-lcc-url').value;

        if (!productId) {
            Utils.showAlert('warning', 'Please select a product first');
            return;
        }

        const btnRegister = document.getElementById('btn-register');
        btnRegister.disabled = true;
        btnRegister.textContent = 'Registering...';

        try {
            const response = await Utils.fetchAPI('/api/instance/register', {
                method: 'POST',
                body: JSON.stringify({
                    product_id: productId,
                    version: version || '1.0.0',
                    lcc_url: lccUrl
                })
            });

            if (response.success) {
                this.instanceData = response;
                this.showRegistrationSuccess(response);
                Utils.showAlert('success', 'Instance registered successfully!');
            } else {
                throw new Error(response.error || 'Registration failed');
            }
        } catch (error) {
            Utils.showAlert('error', `Registration failed: ${error.message}`);
            btnRegister.disabled = false;
            btnRegister.textContent = 'Register Instance';
        }
    },

    showRegistrationSuccess(data) {
        const statusDiv = document.getElementById('registration-status');
        statusDiv.classList.remove('hidden');
        statusDiv.innerHTML = `
            <div style="background: rgba(34, 197, 94, 0.1); border: 1px solid rgba(34, 197, 94, 0.3); border-radius: var(--radius); padding: var(--space-4);">
                <h4 style="color: var(--success); margin-bottom: var(--space-2);">✅ Registration Successful!</h4>
                <div style="color: var(--text-primary);">
                    <div style="margin-bottom: var(--space-1);"><strong>Instance ID:</strong> <code style="color: var(--accent);">${data.instance_id}</code></div>
                    <div style="margin-bottom: var(--space-1);"><strong>Product:</strong> ${data.product_id}</div>
                    <div style="margin-bottom: var(--space-1);"><strong>Version:</strong> ${data.version}</div>
                    <div><strong>Registered:</strong> ${new Date(data.registered_at).toLocaleString()}</div>
                </div>
            </div>
        `;

        const instanceCard = document.getElementById('instance-status-card');
        instanceCard.classList.remove('hidden');

        document.getElementById('instance-id-display').textContent = data.instance_id;
        document.getElementById('product-id-display').textContent = data.product_id;
        document.getElementById('status-display').innerHTML = '<span style="color: var(--success);">✓ Active</span>';
        document.getElementById('registered-at-display').textContent = new Date(data.registered_at).toLocaleString();

        const btnRegister = document.getElementById('btn-register');
        btnRegister.classList.add('hidden');

        const btnClear = document.getElementById('btn-clear-instance');
        btnClear.classList.remove('hidden');
    },

    async clearInstance() {
        const productId = this.selectedProductId;
        if (!productId) {
            Utils.showAlert('warning', 'No instance to clear');
            return;
        }

        try {
            const response = await Utils.fetchAPI('/api/instance/clear', {
                method: 'POST',
                body: JSON.stringify({ product_id: productId })
            });

            if (response.success) {
                document.getElementById('registration-status').classList.add('hidden');
                document.getElementById('instance-status-card').classList.add('hidden');
                
                const btnRegister = document.getElementById('btn-register');
                btnRegister.classList.remove('hidden');
                btnRegister.disabled = false;
                btnRegister.textContent = 'Register Instance';

                const btnClear = document.getElementById('btn-clear-instance');
                btnClear.classList.add('hidden');

                this.instanceData = null;
                Utils.showAlert('success', 'Instance cleared successfully');
            }
        } catch (error) {
            Utils.showAlert('error', `Failed to clear instance: ${error.message}`);
        }
    },

    async testFeature() {
        const productId = this.selectedProductId;
        const featureId = document.getElementById('test-feature-select').value;

        if (!productId || !this.instanceData) {
            Utils.showAlert('warning', 'Please register an instance first');
            return;
        }

        try {
            const response = await Utils.fetchAPI('/api/instance/test', {
                method: 'POST',
                body: JSON.stringify({
                    product_id: productId,
                    feature_id: featureId
                })
            });

            const testResult = document.getElementById('test-result');
            testResult.classList.remove('hidden');

            const resultJson = document.getElementById('test-result-json');
            resultJson.textContent = JSON.stringify(response, null, 2);

            if (response.success && response.enabled) {
                resultJson.style.borderColor = 'var(--success)';
            } else {
                resultJson.style.borderColor = 'var(--error)';
            }
        } catch (error) {
            Utils.showAlert('error', `Test failed: ${error.message}`);
        }
    },

    async generateKeys() {
        try {
            const response = await Utils.fetchAPI('/api/instance/generate-keys', {
                method: 'POST'
            });

            if (response.success) {
                const keysStatus = document.getElementById('keys-status');
                keysStatus.innerHTML = `
                    <div style="color: var(--success); margin-bottom: var(--space-2);">
                        ✓ New key pair generated successfully!
                    </div>
                    <div style="font-size: var(--text-sm); color: var(--text-muted);">
                        Keys are ready for use. They will be used automatically during registration.
                    </div>
                `;
                Utils.showAlert('success', 'Key pair generated successfully');
            }
        } catch (error) {
            Utils.showAlert('error', `Key generation failed: ${error.message}`);
        }
    },

    useSavedKeys() {
        const keysStatus = document.getElementById('keys-status');
        keysStatus.innerHTML = `
            <div style="color: var(--accent);">
                Will attempt to load saved keys from keystore during registration.
                If no keys exist, they will be generated automatically.
            </div>
        `;
        Utils.showAlert('info', 'Will use saved keys if available');
    },

    async viewAllInstances() {
        try {
            const response = await Utils.fetchAPI('/api/instances');
            const instances = response.instances || [];
            
            const panel = document.getElementById('instances-panel');
            const list = document.getElementById('instances-list');
            
            if (instances.length === 0) {
                list.innerHTML = '<div style="color: var(--text-muted); text-align: center; padding: var(--space-4);">No instances registered yet</div>';
            } else {
                list.innerHTML = instances.map(inst => `
                    <div style="background: var(--bg); padding: var(--space-3); border-radius: var(--radius); margin-bottom: var(--space-2); border-left: 3px solid var(--accent);">
                        <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: var(--space-2); color: var(--text-primary);">
                            <div>
                                <div style="font-size: var(--text-sm); color: var(--text-muted);">Instance ID</div>
                                <div style="font-family: var(--font-mono); font-size: var(--text-sm); color: var(--accent); word-break: break-all;">${inst.instance_id.substring(0, 16)}...</div>
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
                                <div style="font-size: var(--text-sm); color: var(--text-muted);">Registered</div>
                                <div style="font-size: var(--text-sm);">${new Date(inst.registered_at).toLocaleString()}</div>
                            </div>
                        </div>
                        <div style="margin-top: var(--space-2); display: flex; gap: var(--space-2);">
                            <button class="btn btn-secondary" style="font-size: var(--text-sm); padding: 4px 8px;" onclick="SetupPage.deleteInstance('${inst.product_id}', '${inst.instance_id}')">Delete</button>
                        </div>
                    </div>
                `).join('');
            }
            
            panel.classList.remove('hidden');
        } catch (error) {
            Utils.showAlert('error', `Failed to load instances: ${error.message}`);
        }
    },

    closeInstancesPanel() {
        document.getElementById('instances-panel').classList.add('hidden');
    },

    async deleteInstance(productId, instanceId) {
        if (!confirm('Are you sure you want to delete this instance?')) return;
        
        try {
            const response = await Utils.fetchAPI('/api/instance/clear', {
                method: 'POST',
                body: JSON.stringify({
                    product_id: productId,
                    instance_id: instanceId
                })
            });

            if (response.success) {
                Utils.showAlert('success', 'Instance deleted successfully');
                await this.viewAllInstances(); // Refresh list
            } else {
                throw new Error(response.error || 'Delete failed');
            }
        } catch (error) {
            Utils.showAlert('error', `Failed to delete instance: ${error.message}`);
        }
    }
};
