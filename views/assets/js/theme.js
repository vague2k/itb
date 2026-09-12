(() => {
	const root = document.documentElement;
	const prefersDark = window.matchMedia("(prefers-color-scheme: dark)");

	const isDark = (mode) => mode === "dark" || (mode === "system" && prefersDark.matches);

	const syncPressed = (mode) => {
		document.querySelectorAll("[data-theme-set]").forEach((el) => {
			el.setAttribute("aria-pressed", String(el.dataset.themeSet === mode));
		});
	};

	const apply = (mode) => {
		root.classList.toggle("dark", isDark(mode));
		root.dataset.theme = mode;
		localStorage.setItem("theme", mode);
		syncPressed(mode);
	};

	document.addEventListener("click", (event) => {
		const button = event.target.closest("[data-theme-set]");
		if (button) apply(button.dataset.themeSet);
	});

	prefersDark.addEventListener("change", () => {
		if ((root.dataset.theme || "system") === "system") apply("system");
	});

	syncPressed(root.dataset.theme || "system");
})();
