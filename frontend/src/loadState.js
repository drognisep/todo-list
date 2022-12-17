export default {
    data: () => {
        return {
            loading: 0,
        }
    },
    methods: {
        startLoading() {
            this.loading++;
        },
        doneLoading() {
            if (this.loading > 0) {
                this.loading--;
            }
        },
    },
    computed: {
        isLoading() {
            return this.loading > 0;
        }
    },
}