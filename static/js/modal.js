window.editModal = function () {
  return {
    show: false,
    init() {
      requestAnimationFrame(() => {
        this.show = true;
        this.$nextTick(() => {
          this.$refs.name?.focus();
        });
      });
    },
    close() {
      this.show = false;
      setTimeout(() => this.$root.remove(), 300);
    },
  };
};

window.layout = function () {
  return {
    atTop: true,

    init() {
      requestAnimationFrame(() => {
        if (window.scrollY > 0) {
          setTimeout(() => {
            this.atTop = false;
          }, 100);
        }
      });
    },

    checkScroll() {
      this.atTop = window.scrollY === 0;
    },
  };
};

window.addEventListener("pageshow", (event) => {
  const nav = performance.getEntriesByType("navigation")[0];
  if (event.persisted || nav?.type === "back_forward") {
    window.location.reload();
  }
});
