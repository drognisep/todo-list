export default {
    setup() {
        window.confirmDialog = function (title, message, onConfirm) {
            let evt = new CustomEvent('confirmDialog', {
                detail: {
                    title, message, onConfirm
                },
            });

            window.dispatchEvent(evt);
        };
    }
}
