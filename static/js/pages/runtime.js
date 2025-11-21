const RuntimePage = {
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container page-enter-active">
                <h1 style="font-size: var(--text-2xl); color: var(--text-secondary);">
                    📊 Live Simulation Dashboard
                </h1>
                <div class="card mt-4">
                    <p>Runtime page - Coming in Week 5</p>
                </div>
            </div>
        `;
    }
};
