import { ref, readonly, watch } from 'vue';

let _state;
export function useReport() {
  if (!_state) {
    const isOpen = ref(false);

    const open = () => { isOpen.value = true; };
    const close = () => { isOpen.value = false; };
    const toggle = () => { isOpen.value = !isOpen.value; };

    // Lock/unlock body scroll when modal is open
    watch(isOpen, (val) => {
      const cls = 'overflow-hidden';
      const el = document?.documentElement || document?.body;
      if (!el) return;
      if (val) el.classList.add(cls);
      else el.classList.remove(cls);
    }, { flush: 'post' });

    _state = {
      isOpen,
      open,
      close,
      toggle,
    };
  }
  return {
    isOpen: readonly(_state.isOpen),
    open: _state.open,
    close: _state.close,
    toggle: _state.toggle,
  };
}
