/**
 * Touch drag-and-drop composable for mobile devices.
 *
 * Mirrors the HTML5 drag API contract (type/data store) using a shared singleton
 * so all draggable elements and drop targets can exchange data exactly like
 * dataTransfer does on desktop — without coupling components directly.
 *
 * Desktop drag events are left untouched; touch events are added in parallel.
 */

const touchDragStore = {
  type: null,
  data: {},
  activeGhost: null,
  dropTarget: null,
}

export function useTouchDragSource(onDragStartCallback, onDragEndCallback) {
  let longPressTimer = null
  let isDragging = false
  let startX = 0
  let startY = 0

  function setDragData(type, payload) {
    touchDragStore.type = type
    touchDragStore.data = payload
  }

  function createGhostElement(sourceElement, touchX, touchY) {
    const rect = sourceElement.getBoundingClientRect()
    const ghost = sourceElement.cloneNode(true)

    ghost.style.cssText = `
      position: fixed;
      top: ${rect.top}px;
      left: ${rect.left}px;
      width: ${rect.width}px;
      height: ${rect.height}px;
      opacity: 0.7;
      pointer-events: none;
      z-index: 9999;
      border-radius: 8px;
      background: #27272a;
      border: 1px solid #3b82f6;
      transition: none;
      box-shadow: 0 8px 24px rgba(0,0,0,0.5);
    `
    document.body.appendChild(ghost)
    touchDragStore.activeGhost = ghost
  }

  function moveGhost(touchX, touchY, originX, originY) {
    if (!touchDragStore.activeGhost) return
    const deltaX = touchX - originX
    const deltaY = touchY - originY
    touchDragStore.activeGhost.style.transform = `translate(${deltaX}px, ${deltaY}px)`
  }

  function removeGhost() {
    if (touchDragStore.activeGhost) {
      touchDragStore.activeGhost.remove()
      touchDragStore.activeGhost = null
    }
  }

  function notifyDropTargetEnter(touchX, touchY) {
    const elementBelow = document.elementFromPoint(touchX, touchY)
    if (!elementBelow) return

    const dropZone = elementBelow.closest('[data-touch-drop-zone]')
    if (dropZone !== touchDragStore.dropTarget) {
      if (touchDragStore.dropTarget) {
        touchDragStore.dropTarget.dispatchEvent(new CustomEvent('touch-drag-leave'))
      }
      touchDragStore.dropTarget = dropZone
      if (dropZone) {
        dropZone.dispatchEvent(new CustomEvent('touch-drag-enter'))
      }
    }
  }

  function handleTouchStart(event, sourceElement, dragDataType, dragDataPayload) {
    if (event.touches.length !== 1) return

    startX = event.touches[0].clientX
    startY = event.touches[0].clientY

    longPressTimer = setTimeout(() => {
      isDragging = true
      setDragData(dragDataType, dragDataPayload)
      createGhostElement(sourceElement, startX, startY)
      onDragStartCallback?.()

      if (navigator.vibrate) navigator.vibrate(40)
    }, 300)
  }

  function handleTouchMove(event) {
    if (!isDragging) {
      clearTimeout(longPressTimer)
      return
    }

    event.preventDefault()

    const touchX = event.touches[0].clientX
    const touchY = event.touches[0].clientY

    moveGhost(touchX, touchY, startX, startY)
    notifyDropTargetEnter(touchX, touchY)
  }

  function handleTouchEnd(event) {
    clearTimeout(longPressTimer)

    if (!isDragging) return

    isDragging = false

    const touchX = event.changedTouches[0].clientX
    const touchY = event.changedTouches[0].clientY

    const elementBelow = document.elementFromPoint(touchX, touchY)
    const dropZone = elementBelow?.closest('[data-touch-drop-zone]')

    if (dropZone) {
      dropZone.dispatchEvent(new CustomEvent('touch-drop', { detail: { ...touchDragStore.data, type: touchDragStore.type } }))
    }

    removeGhost()

    if (touchDragStore.dropTarget) {
      touchDragStore.dropTarget.dispatchEvent(new CustomEvent('touch-drag-leave'))
      touchDragStore.dropTarget = null
    }

    touchDragStore.type = null
    touchDragStore.data = {}

    onDragEndCallback?.()
  }

  return {
    handleTouchStart,
    handleTouchMove,
    handleTouchEnd,
  }
}

export function useTouchDropTarget(onDropCallback, onEnterCallback, onLeaveCallback) {
  function bindDropZone(el) {
    if (!el) return

    el.setAttribute('data-touch-drop-zone', 'true')

    el.addEventListener('touch-drop', (event) => {
      onDropCallback?.(event.detail)
    })

    el.addEventListener('touch-drag-enter', () => {
      onEnterCallback?.()
    })

    el.addEventListener('touch-drag-leave', () => {
      onLeaveCallback?.()
    })
  }

  function unbindDropZone(el) {
    if (!el) return
    el.removeAttribute('data-touch-drop-zone')
  }

  return { bindDropZone, unbindDropZone }
}
