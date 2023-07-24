export const virtualizationTools = {
  updateVisibleElementIndices: (visibleElementIndices, entries, elementIndexFromEntry) => {
    entries.forEach((entry) => {
      const inView = entry.isIntersecting && entry.intersectionRatio >= 0;
      const elementIndex = elementIndexFromEntry(entry);
      if (elementIndex === undefined || elementIndex < 0) {
        return;
      }

      if (inView) {
        visibleElementIndices.add(elementIndex);
      } else {
        /**
         * If the target has no parent-node, it's no longer in the dom-tree.
         * If the element is no longer inView because it was removed from the
         * dom-tree, then this says nothing about the visible indices.
         * If you remove a picture from the grid, the space where the picture
         * was is still visible.
         *
         * We therefore must ignore entries that became invisible that no longer
         * exists
         */
        const entryIsStillMounted = entry.target.parentNode !== null;
        if (entryIsStillMounted) {
          visibleElementIndices.delete(elementIndex);
        }
      }
    });

    /**
     * There are many things that can influence what elements are currently
     * visible on the screen, like scrolling, resizing, menu-opening etc.
     *
     * We therefore cannot make assumptions about our new first- and last
     * visible index, even if it is tempting to initialize these values
     * with this.firstVisibleElementIndex and this.lastVisibleElementIndex.
     *
     * Doing so would break the virtualization though. this.firstVisibleElementIndex
     * would for example always stay at 0
     */
    let firstVisibleElementIndex, lastVisibileElementIndex;
    for (const visibleElementIndex of visibleElementIndices.values()) {
      if (
        firstVisibleElementIndex === undefined ||
        visibleElementIndex < firstVisibleElementIndex
      ) {
        firstVisibleElementIndex = visibleElementIndex;
      }
      if (
        lastVisibileElementIndex === undefined ||
        visibleElementIndex > lastVisibileElementIndex
      ) {
        lastVisibileElementIndex = visibleElementIndex;
      }
    }

    return [firstVisibleElementIndex, lastVisibileElementIndex];
  },

  getVisibleRange: (
    elements,
    containerWidth,
    containerHeight,
    elementWidth,
    elementHeight,
    scrollPos,
    rowPadding = 0
  ) => {
    // floor the elementWidth, so that rounding errors don't cause the columnCount
    // to be one to small
    const visibleColumnCount = Math.floor(containerWidth / Math.floor(elementWidth));

    const totalRowCount = Math.ceil(elements.length / visibleColumnCount);
    const totalScrollHeight = totalRowCount * elementHeight;
    const realScrollPos = Math.max(Math.min(scrollPos, totalScrollHeight), 0);

    /**
     * in a previous version, we calculated the actual visible row count by first
     * calculating the firstRowTopOverlap (= realScrollPos % elementHeight)
     * and adding that to the containerHeight.
     *
     * this would cause the visibleRowCount to sometimes be 1 higher, sometimes
     * 1 smaller than the last run, depending on the actual scroll position.
     *
     * That resulted in this function not always returning the same number of
     * rows. We want to always return the same number of rows however, because
     * that allowes reduces the amount of updates in scenarios where half a row
     * is visible. The result is, that removing and adding a row on scroll always
     * happen in the same render, instead of in two renders
     */
    const visibleRowCount = Math.ceil(containerHeight / elementHeight) + 1;

    const firstVisibleRow = Math.floor(realScrollPos / elementHeight);
    const lastVisibleRow = firstVisibleRow + visibleRowCount - 1;

    const firstRowToRender = firstVisibleRow - rowPadding;
    const lastRowToRender = lastVisibleRow + rowPadding;

    const firstElementToRender = Math.max(firstRowToRender * visibleColumnCount, 0);
    // eslint-disable-next-line prettier-vue/prettier
    const lastElementToRender = Math.min((lastRowToRender + 1) * visibleColumnCount - 1, elements.length - 1);

    return {
      visibleColumnCount,
      firstElementToRender,
      lastElementToRender,
      totalScrollHeight,
    };
  },

  getVirtualizedElementStyle: (index, columnCount, elementWidth, elementHeight) => {
    const rowIndex = Math.floor(index / columnCount);
    const columnIndex = index % columnCount;
    const topPos = rowIndex * elementHeight;
    const leftPos = columnIndex * elementWidth;

    return `display: block; position: absolute; width: ${elementWidth}px; height: ${elementHeight}px; transform: translate(${leftPos}px, ${topPos}px)`;
  },
};
