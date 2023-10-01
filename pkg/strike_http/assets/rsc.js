// @ts-check
import React from "react";
import { jsx, jsxs } from "react/jsx-runtime";
import { StrikeIsland } from "./islands.js";
import { StrikeSuspense } from "./suspense.js";

/** @type {import("./rsc").RscComponent} */
export function RscComponent({ isInitial, url, routerKey, actionData }) {
  if (isInitial) {
    return waitForInitialJSX();
  }
  if (actionData) {
    return fetchClientJSXFromAction(url, routerKey, actionData);
  }
  // return React.use(fetchClientJSX(url, routerKey));
  // Error: Support for `use` not yet implemented in react-debug-tools.
  return fetchClientJSX(url, routerKey);
}

const waitForInitialJSX = React.cache(async function waitForInitialJSX() {
  const chunks = readInitialChunks();
  return await chunksToJSX(chunks);
});

const fetchClientJSX = React.cache(async function fetchClientJSX(href, key) {
  const response = await fetch(href, {
    headers: { RSC: "1" },
  });

  const chunks = readLines(response);
  return await chunksToJSX(chunks);
});

const fetchClientJSXFromAction = React.cache(async function fetchClientJSX(
  href,
  key,
  /** @type {import("./router").ActionData} */
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
  const actionName = "$ACTION_ID_" + actionId;
  formData.append(actionName, "");
  const response = await fetch(href, {
    method: "POST",
    headers: { RSC: "1" },
    body: formData,
  });
  const chunks = readLines(response);
  const ctx = newContext();
  ctx.promises.set(actionName, remotePromise);
  return await chunksToJSX(chunks, ctx);
});

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

export function chunkToJSX(ctx, x) {
  const parsed = JSON.parse(x, function fromJSON(key, value) {
    return parseModelString(ctx, this, key, value);
  });
  // console.log("str", parsed, ctx);
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
  if (
    key === "data-$strike-action" &&
    typeof value === "string" &&
    value.startsWith("$ACTION_ID_")
  ) {
    const actionId = value.substring("$ACTION_ID_".length);
    actionify(parent, actionId);
  }
  if (key === "$strike" && value === "action") {
    actionify(parent, parent.id);
  } else if (key === "$strike" && value === "promise-result") {
    const remote = remotePromiseFromCtx(ctx, parent.id);
    remote.resolve(parent.result);
  } else if (key === "$strike" && value === "promise") {
    const remote = remotePromiseFromCtx(ctx, parent.id);
    promisify(parent, remote.promise);
  } else if (key === "$strike" && value === "component") {
    parent["$$typeof"] = Symbol.for("react.element");
    parent.type = parent["$type"];
    delete parent["$type"];
    parent.ref = null;
    parent.key = parent.props.key || null;
    for (const [k, v] of Object.entries(parent.props || {})) {
      if (k === "style" && typeof v === "string") {
        /** @type {{ [key: string]: string }} */
        const style = {};
        v.split(";").forEach((x) => {
          const [k, v] = x.split(":");
          if (k && v) {
            style[k.trim()] = v.trim();
          }
        });
        parent.props.style = style;
      } else if (k === "key") {
        key = v;
        delete parent.props.key;
      } else if (k === "class") {
        delete parent.props.class;
        parent.props.className = v;
      } else if (parent.type === "meta" && k === "charset") {
        delete parent.props.charset;
        parent.props.charSet = v;
      } else {
        parent.props[k] = v;
      }
    }
    if (parent.type === "strike-suspense") {
      parent.type = StrikeSuspense;
    } else if (parent.type === "strike-island") {
      parent.type = StrikeIsland;
    }
    return undefined;
  }
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
