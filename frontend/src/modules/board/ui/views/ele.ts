// [
//   {
//     x: 569.30859375,
//     y: 267.0234375,
//     id: "qLq8sOgI2OPs0HOQAAS3N",
//     type: "rectangle",
//     width: 164.28515625,
//     height: 193.484375,
//     locked: false,
//     fillStyle: "solid",
//     roughness: 1,
//     roundness: {
//       type: 3,
//     },
//     strokeColor: "#1e1e1e",
//     strokeStyle: "solid",
//     strokeWidth: 2,
//     boundElements: [
//       {
//         id: "JrEDAaTdYA8cReiOXBqT7",
//         type: "arrow",
//       },
//     ],
//   },
//   {
//     x: 912.296875,
//     y: 243.234375,
//     id: "ZDZwK4cakvwjy9kaGecvx",
//     type: "diamond",
//     width: 179.30078125,
//     height: 180.45703125,
//     locked: false,
//     fillStyle: "solid",
//     roughness: 1,
//     roundness: {
//       type: 2,
//     },
//     strokeColor: "#1e1e1e",
//     strokeStyle: "solid",
//     strokeWidth: 2,
//     boundElements: [
//       {
//         id: "JrEDAaTdYA8cReiOXBqT7",
//         type: "arrow",
//       },
//     ],
//   },
//   {
//     x: 651.451171875,
//     y: 363.765625,
//     id: "JrEDAaTdYA8cReiOXBqT7",
//     end: {
//       id: "ZDZwK4cakvwjy9kaGecvx",
//     },
//     type: "arrow",
//     start: {
//       id: "qLq8sOgI2OPs0HOQAAS3N",
//     },
//     width: 350.49609375,
//     height: 30.302734375,
//     locked: false,
//     fillStyle: "solid",
//     roughness: 1,
//     strokeColor: "#1971c2",
//     strokeStyle: "solid",
//     strokeWidth: 2,
//     endArrowhead: "arrow",
//   },
// ];

// const linkElements = useCallback(() => {
//   if (!excalidrawAPI.current) return;

//   const currentElements = excalidrawAPI.current.getSceneElements();

//   console.log("currentElements", currentElements);
//   const nonDeletedElements = currentElements.filter(
//     (el: ExcalidrawElement) => !el.isDeleted && el.type !== "arrow"
//   );

//   if (nonDeletedElements.length < 2) {
//     alert("Need at least 2 elements to create a link");
//     return;
//   }

//   const element1 = nonDeletedElements[nonDeletedElements.length - 2];
//   const element2 = nonDeletedElements[nonDeletedElements.length - 1];

//   const x1 = element1.x + element1.width / 2;
//   const y1 = element1.y + element1.height / 2;
//   const x2 = element2.x + element2.width / 2;
//   const y2 = element2.y + element2.height / 2;

//   const dx = x2 - x1;
//   const dy = y2 - y1;

//   const arrowSkeleton = {
//     type: "arrow" as const,
//     x: x1,
//     y: y1,
//     width: dx,
//     height: dy,
//     strokeColor: "#1971c2",
//     strokeWidth: 2,
//     label: {
//       text: "ARROW",
//       strokeColor: "#099268",
//     },
//     start: {
//       id: element1.id,
//     },
//     end: {
//       id: element2.id,
//     },
//   };

//   const [arrowElement] = convertToExcalidrawElements([arrowSkeleton], {
//     regenerateIds: false,
//   });

//   const arrowId = arrowElement.id || crypto.randomUUID().toString();

//   const arrowWithBindings: any = {
//     ...arrowElement,
//     id: arrowId,
//     startBinding: {
//       elementId: element1.id,
//       focus: 0,
//       gap: 1,
//     },
//     endBinding: {
//       elementId: element2.id,
//       focus: 0,
//       gap: 1,
//     },
//   };

//   const updatedElement1: any = {
//     ...element1,
//     boundElements: [
//       ...(element1.boundElements || []),
//       {
//         id: arrowId,
//         type: "arrow",
//       },
//     ],
//   };

//   const updatedElement2: any = {
//     ...element2,
//     boundElements: [
//       ...(element2.boundElements || []),
//       {
//         id: arrowId,
//         type: "arrow",
//       },
//     ],
//   };

//   // Create updated elements array with modified element1, element2, and new arrow
//   const updatedElements = currentElements.map((el) => {
//     if (el.id === element1.id) {
//       return updatedElement1;
//     }
//     if (el.id === element2.id) {
//       return updatedElement2;
//     }
//     return el;
//   });

//   const finalElements = [...updatedElements, arrowWithBindings];
//   excalidrawAPI.current.updateScene({
//     elements: finalElements,
//   });
// }, []);
