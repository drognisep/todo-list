export default {
    data: () => {
        return {
            loading: 0,
        }
    },
    methods: {
        startLoading() {
            this.waiting++;
        },
        doneLoading() {
            if (this.waiting > 0) {
                this.waiting--;
            }
        },
    },
    computed: {
        isLoading() {
            return this.loading > 0;
        }
    },
}