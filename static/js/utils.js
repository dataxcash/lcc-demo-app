const Utils = {
    async fetchAPI(url, options = {}) {
        try {
            const response = await fetch(url, {
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                },
                ...options
            });
            
            if (!response.ok) {
                const error = await response.json().catch(() => ({}));
                throw new Error(error.message || `HTTP ${response.status}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error(`API Error [${url}]:`, error);
            throw error;
        }
    },

    showAlert(type, message) {
        const existingAlert = document.querySelector('.alert');
        if (existingAlert) {
            existingAlert.remove();
        }

        const alert = document.createElement('div');
        alert.className = `alert alert-${type}`;
        alert.textContent = message;
        
        const container = document.querySelector('.page-container');
        if (container) {
            container.insertBefore(alert, container.firstChild);
            
            setTimeout(() => {
                alert.style.opacity = '0';
                setTimeout(() => alert.remove(), 300);
            }, 5000);
        }
    },

    updateLCCStatus(connected, details = {}) {
        const statusElement = document.querySelector('.lcc-status');
        const statusText = document.querySelector('.status-text');
        
        if (connected) {
            statusElement.classList.remove('disconnected');
            statusElement.classList.add('connected');
            statusText.textContent = details.version 
                ? `Connected | ${details.version}`
                : 'Connected';
        } else {
            statusElement.classList.remove('connected');
            statusElement.classList.add('disconnected');
            statusText.textContent = 'Disconnected';
        }
    },

    setActiveStep(pageName) {
        document.querySelectorAll('.step').forEach(step => {
            if (step.dataset.page === pageName) {
                step.classList.add('active');
            } else {
                step.classList.remove('active');
            }
        });
    },

    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    },

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
};
