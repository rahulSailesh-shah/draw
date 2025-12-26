[
  {
    id: "rectangle-1766697651637-bnbxm8giz",
    type: "rectangle",
    x: -297.35935639756735,
    y: -218.5405722658417,
    width: 150,
    height: 100,
    angle: 0,
    strokeColor: "#1e1e1e",
    backgroundColor: "#c0eb75",
    fillStyle: "solid",
    strokeWidth: 2,
    strokeStyle: "solid",
    roughness: 1,
    opacity: 100,
    groupIds: [],
    frameId: null,
    index: "a0",
    roundness: null,
    seed: 292987561,
    version: 83,
    versionNonce: 1995346537,
    isDeleted: false,
    boundElements: null,
    updated: 1766697654224,
    link: null,
    locked: false,
  },
  {
    id: "ellipse-1766697655654-8gzkuat8u",
    type: "ellipse",
    x: -84.4794756447319,
    y: -216.36837463270717,
    width: 150,
    height: 100,
    angle: 0,
    strokeColor: "#1e1e1e",
    backgroundColor: "#ffc9c9",
    fillStyle: "solid",
    strokeWidth: 2,
    strokeStyle: "solid",
    roughness: 1,
    opacity: 100,
    groupIds: [],
    frameId: null,
    index: "a1",
    roundness: null,
    seed: 968079303,
    version: 242,
    versionNonce: 665886249,
    isDeleted: false,
    boundElements: null,
    updated: 1766697661458,
    link: null,
    locked: false,
  },
  {
    id: "uHmzgOxa6_-LTeztfGO7W",
    type: "arrow",
    x: -222.35935639756735,
    y: -168.5405722658417,
    width: 212.87988075283545,
    height: 2.172197633134516,
    angle: 0,
    strokeColor: "#1971c2",
    backgroundColor: "transparent",
    fillStyle: "solid",
    strokeWidth: 2,
    strokeStyle: "solid",
    roughness: 1,
    opacity: 100,
    groupIds: [],
    frameId: null,
    index: "a2",
    roundness: null,
    seed: 404570119,
    version: 3,
    versionNonce: 1856349991,
    isDeleted: false,
    boundElements: null,
    updated: 1766697662671,
    link: null,
    locked: false,
    points: [
      [0.5, 0.5],
      [212.37988075283545, 1.672197633134516],
    ],
    lastCommittedPoint: null,
    startBinding: {
      elementId: "rectangle-1766697651637-bnbxm8giz",
      focus: 0,
      gap: 1,
    },
    endBinding: {
      elementId: "ellipse-1766697655654-8gzkuat8u",
      focus: 0,
      gap: 1,
    },
    startArrowhead: null,
    endArrowhead: "arrow",
    elbowed: false,
  },
  {
    id: "UYbpIqbZkSPjoWzVQ-St0",
    type: "diamond",
    x: 391.0370708106983,
    y: 438.4103121144316,
    width: 136.015625,
    height: 169.513916015625,
    angle: 0,
    strokeColor: "#1e1e1e",
    backgroundColor: "transparent",
    fillStyle: "solid",
    strokeWidth: 2,
    strokeStyle: "solid",
    roughness: 1,
    opacity: 100,
    groupIds: [],
    frameId: null,
    index: "a3",
    roundness: {
      type: 2,
    },
    seed: 1854228839,
    version: 63,
    versionNonce: 486685577,
    isDeleted: false,
    boundElements: [
      {
        id: "Dpb8idyI6ziQsrzIFXYGf",
        type: "arrow",
      },
    ],
    updated: 1766697690906,
    link: null,
    locked: false,
  },
  {
    id: "3z4_7q-cewxqM59GK1xxR",
    type: "diamond",
    x: 667.0309672950733,
    y: 443.9789156300566,
    width: 136.8707275390625,
    height: 162.06597900390625,
    angle: 0,
    strokeColor: "#1e1e1e",
    backgroundColor: "transparent",
    fillStyle: "solid",
    strokeWidth: 2,
    strokeStyle: "solid",
    roughness: 1,
    opacity: 100,
    groupIds: [],
    frameId: null,
    index: "a4",
    roundness: {
      type: 2,
    },
    seed: 801259657,
    version: 287,
    versionNonce: 1692995623,
    isDeleted: false,
    boundElements: [
      {
        id: "Dpb8idyI6ziQsrzIFXYGf",
        type: "arrow",
      },
    ],
    updated: 1766697697600,
    link: null,
    locked: false,
  },
  {
    id: "Dpb8idyI6ziQsrzIFXYGf",
    type: "arrow",
    x: 537.3129023624884,
    y: 522.7980922671292,
    width: 122.1998817173187,
    height: 2.206473280032469,
    angle: 0,
    strokeColor: "#1e1e1e",
    backgroundColor: "transparent",
    fillStyle: "solid",
    strokeWidth: 2,
    strokeStyle: "solid",
    roughness: 1,
    opacity: 100,
    groupIds: [],
    frameId: null,
    index: "a5",
    roundness: {
      type: 2,
    },
    seed: 2002086697,
    version: 291,
    versionNonce: 767718215,
    isDeleted: false,
    boundElements: null,
    updated: 1766697697601,
    link: null,
    locked: false,
    points: [
      [0, 0],
      [122.1998817173187, 2.206473280032469],
    ],
    lastCommittedPoint: null,
    startBinding: {
      elementId: "UYbpIqbZkSPjoWzVQ-St0",
      focus: -0.019315047302223245,
      gap: 14.595256120282956,
    },
    endBinding: {
      elementId: "3z4_7q-cewxqM59GK1xxR",
      focus: -0.017466499813780206,
      gap: 11.903051894184799,
    },
    startArrowhead: null,
    endArrowhead: "arrow",
    elbowed: false,
  },
];

const addDiamond = useCallback(() => {
  if (!excalidrawAPI.current) return;

  const currentElements = excalidrawAPI.current.getSceneElements();
  const newElement = convertToExcalidrawElements(
    [
      {
        type: "diamond",
        id: generateElementId("diamond"),
        x: 100 + Math.random() * 200,
        y: 100 + Math.random() * 200,
        width: 150,
        height: 100,
        backgroundColor: "#a5d8ff",
        strokeWidth: 2,
      },
    ],
    { regenerateIds: false }
  )[0];

  const updatedElements = [...currentElements, newElement];
  excalidrawAPI.current.updateScene({
    elements: updatedElements,
  });

  excalidrawAPI.current.scrollToContent([newElement]);
}, []);

const linkElements = useCallback(() => {
  if (!excalidrawAPI.current) return;

  const currentElements = excalidrawAPI.current.getSceneElements();

  console.log(currentElements);
  const nonDeletedElements = currentElements.filter(
    (el: ExcalidrawElement) => !el.isDeleted && el.type !== "arrow"
  );

  if (nonDeletedElements.length < 2) {
    alert("Need at least 2 elements to create a link");
    return;
  }

  const element1 = nonDeletedElements[nonDeletedElements.length - 2];
  const element2 = nonDeletedElements[nonDeletedElements.length - 1];

  const x1 = element1.x + element1.width / 2;
  const y1 = element1.y + element1.height / 2;
  const x2 = element2.x + element2.width / 2;
  const y2 = element2.y + element2.height / 2;

  const dx = x2 - x1;
  const dy = y2 - y1;

  const arrowSkeleton = {
    type: "arrow" as const,
    x: x1,
    y: y1,
    width: dx,
    height: dy,
    strokeColor: "#1971c2",
    strokeWidth: 2,
    start: {
      id: element1.id,
    },
    end: {
      id: element2.id,
    },
  };

  const [arrowElement] = convertToExcalidrawElements([arrowSkeleton], {
    regenerateIds: false,
  });

  const arrowId = arrowElement.id || generateElementId("arrow");

  const arrowWithBindings: any = {
    ...arrowElement,
    id: arrowId,
    startBinding: {
      elementId: element1.id,
      focus: 0,
      gap: 1,
    },
    endBinding: {
      elementId: element2.id,
      focus: 0,
      gap: 1,
    },
  };

  const updatedElement1: any = {
    ...element1,
    boundElements: [
      ...(element1.boundElements || []),
      {
        id: arrowId,
        type: "arrow",
      },
    ],
  };

  const updatedElement2: any = {
    ...element2,
    boundElements: [
      ...(element2.boundElements || []),
      {
        id: arrowId,
        type: "arrow",
      },
    ],
  };

  // Create updated elements array with modified element1, element2, and new arrow
  const updatedElements = currentElements.map((el) => {
    if (el.id === element1.id) {
      return updatedElement1;
    }
    if (el.id === element2.id) {
      return updatedElement2;
    }
    return el;
  });

  const finalElements = [...updatedElements, arrowWithBindings];
  excalidrawAPI.current.updateScene({
    elements: finalElements,
  });
}, []);

const generateElementId = (type: string) => {
  return `${type}-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
};
