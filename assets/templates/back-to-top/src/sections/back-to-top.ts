const btn = document.getElementById("backToTop");

if (btn) {
  function handleShowButton() {
    if (!btn) return;
    if (window.scrollY > 100) {
      btn.setAttribute("data-visible", "true");
      btn.removeAttribute("tabindex");
    } else {
      btn.removeAttribute("data-visible");
      btn.setAttribute("tabindex", "-1");
    }
  }

  window.addEventListener("scroll", handleShowButton);

  window.addEventListener("DOMContentLoaded", handleShowButton);

  btn.addEventListener("click", () => {
    const prefersReducedMotion = window.matchMedia?.(
      "(prefers-reduced-motion: reduce)",
    ).matches;
    const root = document.scrollingElement || document.documentElement;

    const scrollOptions = prefersReducedMotion
      ? ({ top: 0 } satisfies ScrollToOptions)
      : ({ top: 0, behavior: "smooth" } satisfies ScrollToOptions);

    root.scrollTo(scrollOptions);

    const main = document.querySelector("main");
    if (main) {
      main.setAttribute("tabindex", "-1");
      main.focus({ preventScroll: true });
    }
  });
}
