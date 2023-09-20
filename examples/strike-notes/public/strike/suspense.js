import React from "react";
import { jsx } from "react/jsx-runtime";

export function StrikeSuspense(props) {
  console.log("StrikeSuspense", props.children[0]);
  const [isMounted, setIsMounted] = React.useState(false);
  React.useEffect(() => {
    setIsMounted(true);
  }, []);
  if (props.cantStream) {
    return props.children;
  }
  if (!isMounted) {
    // TODO: fix hydration
    return props.fallback;
  }
  return jsx(React.Suspense, {
    fallback: props.fallback,
    children: jsx(Render, { children: props.children }),
  });
}

function Render(props) {
  return props.children;
}
