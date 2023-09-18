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
  console.log("router", router);
  React.useEffect(() => {
    addNavigation(setRouter);
  }, []);
  return jsx(RscComponent, {
    url: router.path,
    routerKey: router.key,
  });
}

/** @type {import("./router").createRouterState} */
function createRouterState(path) {
  // compare this to https://github.com/vercel/next.js/blob/c6c38916882e419d9c4babdd9223339094fff1c3/packages/next-swc/crates/next-core/js/src/entry/app/hydrate.tsx#L130

  // if (typeof __rsc === "undefined") {
  //   window.__rsc = {
  //     push: function (x) {
  //       window.__rsc = [x];
  //       boot();
  //     },
  //   };
  // } else {
  //   boot();
  // }

  return { path, isInitial: true, key: "initial" };
}

/** @type {import("./router").changeRouterState} */
function changeRouterState(path, key) {
  return { path, isInitial: false, key };
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
