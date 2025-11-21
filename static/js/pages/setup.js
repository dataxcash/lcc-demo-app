const SetupPage = {
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container page-enter-active">
                <h1 style="font-size: var(--text-2xl); color: var(--text-secondary);">
                    ⚙️ Configure Simulation Instance
                </h1>
                <div class="card mt-4">
                    <p>Setup page - Coming in Week 4</p>
                </div>
            </div>
        `;
    }
};
