export default {
    setup() {
        window.showProgress = function (message) {
            let evt = new CustomEvent('showProgressDialog', {
                detail: {
                    message,
                }
            })

            window.dispatchEvent(evt);
        };

        window.closeProgress = function () {
            let evt = new CustomEvent('closeProgressDialog')
            window.dispatchEvent(evt);
        };
    }
}