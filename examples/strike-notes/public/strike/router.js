// @ts-check
import { RscComponent, jsonToJSX } from "./rsc.js";
import React from "react";
import { jsx } from "react/jsx-runtime";

/** @type {import("./react").useState} */
const useState = React.useState;

/** @type {import("./router").Router} */
export function Router() {
  const [router, setRouter] = useState(() =>
    createRouterState(window.location.pathname)
  );
  React.useEffect(() => {
    addNavigation(setRouter);
  }, []);
  return jsx(RscComponent, {
    initialPage: router.initialPage,
    url: router.path,
    routerKey: router.key,
  });
}

/** @type {import("./router").createRouterState} */
function createRouterState(path) {
  const jsonData = JSON.parse(window.__rsc[0]);
  const page = jsonToJSX(jsonData);
  return { path, isInitial: true, initialPage: page, key: "initial" };
}

/** @type {import("./router").changeRouterState} */
function changeRouterState(path, key) {
  return { path, isInitial: false, key, initialPage: null };
}

/** @type {import("./router").addNavigation} */
function addNavigation(setRouter) {
  /** @type {import("./router").navigate} */
  function navigate(pathname) {
    React.startTransition(() => {
      // invalidate the cache on every navigation
      setRouter(changeRouterState(pathname, "" + Math.random()));
    });
  }

  window.addEventListener(
    "click",
    (e) => {
      if (e.target.tagName !== "A") {
        return;
      }
      if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) {
        return;
      }
      const href = e.target.getAttribute("href");
      if (!href.startsWith("/")) {
        return;
      }
      e.preventDefault();
      window.history.pushState(null, "", href);
      navigate(href);
    },
    true
  );

  window.addEventListener("popstate", () => {
    navigate(window.location.pathname);
  });
  window.__rscNav = (pathname) => {
    window.history.pushState(null, "", pathname);
    navigate(pathname);
  };
}
