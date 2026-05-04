document.addEventListener('DOMContentLoaded', () => {
    function setupTabs(groupSelector, btnSelector, contentSelector) {
        const groups = document.querySelectorAll(groupSelector);
        
        groups.forEach(group => {
            const buttons = group.querySelectorAll(btnSelector);
            const contents = group.parentElement.querySelectorAll(contentSelector);

            buttons.forEach(btn => {
                btn.addEventListener('click', () => {
                    const tabId = btn.dataset.tab;

                    // Update buttons in this group
                    buttons.forEach(b => b.classList.remove('active'));
                    btn.classList.add('active');

                    // Update contents for this group
                    contents.forEach(content => {
                        if (content.parentElement === group.parentElement) {
                            content.classList.remove('active');
                            if (content.id === tabId || content.dataset.subTab === tabId) {
                                content.classList.add('active');
                            }
                        }
                    });
                });
            });
        });
    }

    // Setup main tabs
    const mainButtons = document.querySelectorAll('.tab-btn');
    const mainContents = document.querySelectorAll('.tab-content');
    mainButtons.forEach(btn => {
        btn.addEventListener('click', () => {
            mainButtons.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            mainContents.forEach(c => {
                c.classList.remove('active');
                if (c.id === btn.dataset.tab) c.classList.add('active');
            });
        });
    });

    // Setup sub-tabs
    document.addEventListener('click', (e) => {
        if (e.target.classList.contains('sub-tab-btn')) {
            const btn = e.target;
            const group = btn.parentElement;
            const tabId = btn.dataset.tab;
            
            group.querySelectorAll('.sub-tab-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            
            group.parentElement.querySelectorAll('.sub-tab-content').forEach(c => {
                c.classList.remove('active');
                if (c.dataset.id === tabId) c.classList.add('active');
            });
        }
    });
});
