const App = {
    state: {
        currentPage: 'welcome',
        config: null
    },

    pages: {
        welcome: WelcomePage,
        tiers: TiersPage,
        limits: LimitsPage,
        setup: SetupPage,
        runtime: RuntimePage
    },

    pageOrder: ['welcome', 'tiers', 'limits', 'setup', 'runtime'],

    async init() {
        this.setupRouter();
        this.setupNavigation();
        await this.renderCurrentPage();
    },

    setupRouter() {
        window.addEventListener('hashchange', () => {
            this.renderCurrentPage();
        });

        window.addEventListener('popstate', () => {
            this.renderCurrentPage();
        });
    },

    setupNavigation() {
        const steps = document.querySelectorAll('.step');
        steps.forEach(step => {
            step.addEventListener('click', () => {
                const page = step.dataset.page;
                this.navigateTo(page);
            });
        });

        const btnBack = document.getElementById('btn-back');
        const btnNext = document.getElementById('btn-next');

        btnBack.addEventListener('click', () => this.goBack());
        btnNext.addEventListener('click', () => this.goNext());
    },

    async renderCurrentPage() {
        const hash = window.location.hash.slice(1);
        const page = hash || 'welcome';
        
        if (!this.pages[page]) {
            this.navigateTo('welcome');
            return;
        }

        this.state.currentPage = page;
        
        Utils.setActiveStep(page);
        
        this.updateNavButtons();
        
        const pageModule = this.pages[page];
        if (pageModule && pageModule.render) {
            await pageModule.render();
        }
    },

    navigateTo(page) {
        if (this.pages[page]) {
            window.location.hash = page;
        }
    },

    goBack() {
        const currentIndex = this.pageOrder.indexOf(this.state.currentPage);
        if (currentIndex > 0) {
            this.navigateTo(this.pageOrder[currentIndex - 1]);
        }
    },

    goNext() {
        const currentIndex = this.pageOrder.indexOf(this.state.currentPage);
        if (currentIndex < this.pageOrder.length - 1) {
            this.navigateTo(this.pageOrder[currentIndex + 1]);
        }
    },

    updateNavButtons() {
        const btnBack = document.getElementById('btn-back');
        const btnNext = document.getElementById('btn-next');
        const currentIndex = this.pageOrder.indexOf(this.state.currentPage);

        btnBack.disabled = currentIndex === 0;
        btnNext.disabled = currentIndex === this.pageOrder.length - 1;

        if (currentIndex === 0) {
            btnBack.style.visibility = 'hidden';
        } else {
            btnBack.style.visibility = 'visible';
        }

        if (currentIndex === this.pageOrder.length - 1) {
            btnNext.style.visibility = 'hidden';
        } else {
            btnNext.style.visibility = 'visible';
        }
    }
};

document.addEventListener('DOMContentLoaded', () => {
    App.init();
});
