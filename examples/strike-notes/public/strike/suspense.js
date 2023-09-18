import React from "react";

export function StrikeSuspense(props) {
  console.log("StrikeSuspense", props.children[0]);
  const [isMounted, setIsMounted] = React.useState(false);
  React.useEffect(() => {
    setIsMounted(true);
  }, []);
  if (isMounted) {
    return props.children;
  }
  return props.fallback;
}
