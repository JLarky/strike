// @ts-check
import {
  RscComponent,
  createRemotePromise,
  fetchChunksPromise,
  fetchFromActionPromise,
} from "./rsc.js";
import React from "react";
import { jsx } from "react/jsx-runtime";

/** @type {import("./react").useState} */
const useState = React.useState;

/** @type {import("./router.js").Router} */
export function Router() {
  const [router, setRouter] = useState(() =>
    createRouterState(window.location.pathname + window.location.search)
  );
  React.useEffect(() => {
    addNavigation(setRouter);
  }, []);
  return jsx(
    RscComponent,
    /** @satisfies {import("./rsc").RscComponentProps} */ ({
      isInitial: router.isInitial,
      url: router.href,
      urlPromise: router.urlPromise,
      routerKey: router.key,
      actionData: router.actionData,
      actionPromise: router.actionPromise,
    })
  );
}

/** @type {import("./router.js").createRouterState} */
function createRouterState(href) {
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

  return {
    href,
    urlPromise: undefined,
    actionPromise: undefined,
    isInitial: true,
    key: "initial",
  };
}

/** @type {import("./router.js").changeRouterState} */
function changeRouterState(href, key) {
  const urlPromise = fetchChunksPromise(href);
  return {
    href,
    urlPromise,
    actionPromise: undefined,
    isInitial: false,
    key,
  };
}

/** @type {import("./router.js").changeRouterStateForAction} */
function changeRouterStateForAction(href, key, actionData) {
  const actionPromise = fetchFromActionPromise(href, actionData);
  return {
    href,
    urlPromise: undefined,
    actionPromise,
    isInitial: false,
    key,
    actionData,
  };
}

/** @type {import("./router.js").addNavigation} */
function addNavigation(setRouter) {
  /** @type {import("./router.js").navigate} */
  function navigate(href) {
    React.startTransition(() => {
      // invalidate the cache on every navigation
      setRouter(changeRouterState(href, "" + Math.random()));
    });
  }
  /** @type {import("./router.js").submitForm} */
  function submitForm(actionData) {
    React.startTransition(() => {
      // invalidate the cache on every navigation
      setRouter((x) =>
        changeRouterStateForAction(x.href, "" + Math.random(), actionData)
      );
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
    navigate(window.location.pathname + window.location.search);
  });
  window.__rscNav = (href) => {
    window.history.pushState(null, "", href);
    navigate(href);
  };
  /** @type {typeof window.__rscAction} */
  window.__rscAction = (actionId0, data) => {
    let actionId = actionId0;
    if (!actionId.startsWith("$ACTION_ID_")) {
      actionId = "$ACTION_ID_" + actionId;
    }
    const remotePromise = createRemotePromise(actionId);
    submitForm({ actionId, data, remotePromise });
    return remotePromise.promise;
  };
}
