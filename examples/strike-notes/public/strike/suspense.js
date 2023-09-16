import React from "react";
import { jsx, jsxs } from "react/jsx-runtime";

export function StrikeSuspense(props) {
  console.log("StrikeSuspense", props);
  const [isMounted, setIsMounted] = React.useState(false);
  React.useEffect(() => {
    setIsMounted(true);
  }, []);
  if (isMounted && false) {
    return jsx(React.Suspense, { fallback: "" }, "");
  }
  // TODO: seems wrong, but I get hydration errors if I don't do this
  return props.fallback;
}
