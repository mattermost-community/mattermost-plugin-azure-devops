import {useEffect} from 'react';

/**
 * Hook that detects clicks outside of the passed ref
 */
function useOutsideClick(ref: React.RefObject<HTMLInputElement>, handleOnClickOutside: () => void) {
    useEffect(() => {
        /**
         * Check if clicked on outside of element
         */
        function handleClickOutside(event: MouseEvent) {
            if (ref.current && !ref.current.contains(event.target as unknown as Node)) {
                handleOnClickOutside();
            }
        }

        // Bind the event listener
        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            // Unbind the event listener on clean up
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [ref]);
}

export default useOutsideClick;
