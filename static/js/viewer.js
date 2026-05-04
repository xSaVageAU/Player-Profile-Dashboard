document.addEventListener('DOMContentLoaded', () => {
    const viewerEl = document.getElementById('skin-viewer');
    if (!viewerEl) return;

    const skinUrl = viewerEl.dataset.skin || 'https://mineskin.org/render/head?url=http://textures.minecraft.net/texture/3ff95f4e1f7535b917535b917535b917535b917535b917535b917535b917535b';

    if (typeof skinview3d === 'undefined') {
        console.error('SkinView3D library not loaded. Check your internet connection or CDN link.');
        viewerEl.innerHTML = '<div style="display:flex; align-items:center; justify-content:center; height:100%; border:1px dashed var(--border-color); border-radius:24px; color:var(--text-secondary)">Library Load Failed</div>';
        return;
    }

    setTimeout(() => {
        try {
            // Ensure dimensions are valid
            const width = viewerEl.offsetWidth || 300;
            const height = viewerEl.offsetHeight || 500;

            const skinViewer = new skinview3d.SkinViewer({
                canvas: document.createElement('canvas'),
                width: width,
                height: height,
                skin: skinUrl
            });

            viewerEl.innerHTML = ''; // Clear fallback text
            viewerEl.appendChild(skinViewer.canvas);

            // Use built-in controls (v3.x+)
            skinViewer.controls.enableRotate = true;
            skinViewer.controls.enableZoom = false; 
            skinViewer.controls.enablePan = false;

            // Disable automatic movement
            skinViewer.autoRotate = false;

            // Responsive resize
            window.addEventListener('resize', () => {
                skinViewer.width = viewerEl.offsetWidth;
                skinViewer.height = viewerEl.offsetHeight;
            });
        } catch (e) {
            console.error('Failed to load SkinView3D:', e);
            viewerEl.innerHTML = '<div style="display:flex; align-items:center; justify-content:center; height:100%; border:1px dashed var(--border-color); border-radius:24px; color:var(--text-secondary)">3D Preview Unavailable</div>';
        }
    }, 100);
});
