// @ts-check
import React from "react";
import { jsx, jsxs } from "react/jsx-runtime";
import { StrikeIsland } from "./islands.js";
import { StrikeSuspense } from "./suspense.js";

const __debug = {
  /** @type {any[]} */
  chunks: [],
  /** @type {any[]} */
  jsx: [],
};
// @ts-ignore
window.__debug = __debug;

/** @type {import("./rsc").RscComponent} */
export function RscComponent({
  isInitial,
  url,
  urlPromise,
  routerKey,
  actionPromise,
  actionData,
}) {
  if (isInitial) {
    return waitForInitialJSX();
  }
  if (actionData) {
    return fetchClientJSXFromAction(actionPromise);
  }
  return fetchClientJSX(urlPromise);
}

const initialChunksPromise = (() => {
  const chunks = readInitialChunks();
  return chunksToJSX(chunks);
})();

function waitForInitialJSX() {
  return React.use(initialChunksPromise);
}

/** @type {import("./rsc").fetchChunksPromise} */
export const fetchChunksPromise = async (href) => {
  const response = await fetch(href, {
    headers: { RSC: "1" },
  });

  const chunks = readLines(response);
  return await chunksToJSX(chunks);
};

function fetchClientJSX(urlPromise) {
  return React.use(urlPromise);
}

/** @type {import("./rsc").fetchFromActionPromise} */
export const fetchFromActionPromise = async function fetchFromActionPromise(
  href,
  actionData
) {
  const { actionId, data, remotePromise } = actionData;
  // convert everything to FormData
  /** @type {FormData | undefined} */
  let formData = data instanceof FormData ? data : undefined;
  if (!formData) {
    formData = new FormData();
    if (data !== undefined) {
      formData.append("data", JSON.stringify(data));
    }
  }
  for (const k of formData.keys()) {
    if (k.startsWith("$ACTION_ID_")) {
      formData.delete(k);
    }
  }
  const actionName = actionId;
  formData.append(actionName, "");
  const response = await fetch(href, {
    method: "POST",
    headers: { RSC: "1" },
    body: formData,
  }).catch((e) => {
    remotePromise.reject(e);
  });
  const chunks = readLines(response);
  const ctx = newContext();
  // FIXME: ctx.promises.set(actionName, remotePromise);
  const jsx = await chunksToJSX(chunks, ctx).catch((e) => {
    remotePromise.reject(e);
  });
  remotePromise.resolve("done");
  return jsx;
};

function fetchClientJSXFromAction(actionPromise) {
  return React.use(actionPromise);
}

function newContext() {
  return { promises: new Map() };
}

async function chunksToJSX(chunks, ctx = newContext()) {
  const root = await chunks.next().then((x) => chunkToJSX(ctx, x.value));
  (async () => {
    for await (const line of chunks) {
      chunkToJSX(ctx, line);
    }
  })();
  return root;
}

/** @type {import("./rsc").chunkToJSX}*/
export function chunkToJSX(ctx, x) {
  __debug.chunks.push(JSON.parse(x));
  const parsed = JSON.parse(x, function fromJSON(key, value) {
    return parseModelString(ctx, this, key, value);
  });
  // console.log("str", parsed, ctx);
  __debug.jsx.push(parsed);
  return parsed;
}

/** @type {import("./rsc").createRemotePromise}*/
export function createRemotePromise(id) {
  /** @type {(value: any) => void} */
  let resolve = () => {};
  let reject = () => {};
  const promise = new Promise((res, rej) => {
    resolve = res;
    reject = rej;
  });
  return { id, promise, resolve, reject };
}

/** @type {import("./rsc").remotePromiseFromCtx}*/
function remotePromiseFromCtx(ctx, id) {
  let remote = ctx.promises.get(id);
  if (!remote) {
    remote = createRemotePromise(id);
    ctx.promises.set(id, remote);
  }
  return remote;
}

/** @type {import("./rsc").promisify}*/
function promisify(obj, promise) {
  obj.__proto__ = promise.__proto__;
  obj.promise = promise;
  obj.then = promise.then.bind(promise);
  obj.catch = promise.catch.bind(promise);
  obj.finally = promise.finally.bind(promise);
}

/** @type {import("./rsc").actionify}*/
function actionify(obj, actionId) {
  obj.action = function (formData) {
    return window.__rscAction(actionId, formData);
  };
}

/** @type {import("./rsc").parseModelString} */
function parseModelString(ctx, parent, key, value) {
  if (Array.isArray(value)) {
    if (value[0] === "$strike:element") {
      const { key, ...props } = value[2];
      return jsxs(value[1], props, key);
    } else if (value[0] === "$strike:text") {
      return value[1];
    } else if (value[0] === "$strike:island") {
      return jsx(StrikeIsland, {
        component: value[1],
        islandProps: value[2],
        ssrFallback: value[3],
      });
    } else if (value[0] === "$strike:island-go") {
      const {
        "component-export": component,
        ssrFallback,
        ...islandProps
      } = value[1];
      return jsxs(StrikeIsland, {
        component,
        islandProps,
        ssrFallback,
      });
    } else if (value[0] === "$strike:form") {
      // fixes `Cannot specify a encType or method for a form that specifies a function as the action. React provides those automatically. They will get overridden.`
      const {
        key,
        encType,
        method,
        ["data-$strike-action"]: id,
        ...props
      } = value[1];
      actionify(props, id);
      return jsxs("form", props, key);
    }
  }
  // if (key === "$strike" && value === "action") {
  //   actionify(parent, parent.id);
  // } else if (key === "$strike" && value === "promise-result") {
  //   const remote = remotePromiseFromCtx(ctx, parent.id);
  //   remote.resolve(parent.result);
  // } else if (key === "$strike" && value === "promise") {
  //   const remote = remotePromiseFromCtx(ctx, parent.id);
  //   promisify(parent, remote.promise);
  // }
  return value;
}

async function* readInitialChunks() {
  // make sure that it's an array (this code could be executed before JSX streaming has started)
  window.__rsc = window.__rsc || [];
  // consume anything that was already streamed
  for (const x of __rsc) {
    yield x;
  }
  // wait for new chunks to be streamed
  const chunkQueue = [];
  let resolveChunk = null;

  // new chunks will call __rsc.push
  const originalPush = window.__rsc.push.bind(window.__rsc);
  window.__rsc.push = function (x) {
    originalPush(x); // to make array look pretty
    chunkQueue.push(x); // collect chunks
    if (resolveChunk) {
      resolveChunk(); // trigger promise
      resolveChunk = null;
    }
  };

  // keep generator running if we are still waiting for chunks
  while (true) {
    if (chunkQueue.length > 0) {
      yield chunkQueue.shift();
    } else {
      // Wait for the next chunk
      // TODO: stop stream when document is fully loaded
      await new Promise((resolve) => {
        resolveChunk = resolve;
      });
    }
  }
}

async function* readLines(response) {
  const reader = response.body?.getReader();
  let accumulatedData = "";

  if (reader) {
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;

      // Convert Uint8Array to a string
      accumulatedData += new TextDecoder().decode(value);

      // Split by line breaks but keep the last, potentially incomplete line
      let lastNewlineIndex = accumulatedData.lastIndexOf("\n\n");
      if (lastNewlineIndex !== -1) {
        const lines = accumulatedData
          .substring(0, lastNewlineIndex)
          .split("\n\n");
        for (const line of lines) {
          yield line;
        }

        // Keep the remainder for the next iteration
        accumulatedData = accumulatedData.substring(lastNewlineIndex + 2);
      }
    }

    // If there's any remaining data after all chunks have been processed, yield it
    if (accumulatedData) {
      yield accumulatedData;
    }
  }
}
