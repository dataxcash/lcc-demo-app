const TiersPage = {
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container page-enter-active">
                <h1 style="font-size: var(--text-2xl); color: var(--text-secondary);">
                    📚 Lesson 1: Understanding License Tiers
                </h1>
                <div class="card mt-4">
                    <p>Tiers page - Coming in Week 2</p>
                </div>
            </div>
        `;
    }
};
