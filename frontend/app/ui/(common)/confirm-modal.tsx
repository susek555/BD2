import * as AlertDialog from '@radix-ui/react-alert-dialog';

interface ConfirmationModalProps {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  onConfirm: () => void;
  onCancel: () => void;
  isOpen: boolean;
}

export default function ConfirmationModal({
  title,
  message,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  onConfirm,
  onCancel,
  isOpen,
}: ConfirmationModalProps) {
  return (
    <AlertDialog.Root
      open={isOpen}
      onOpenChange={(open) => {
        if (!open) onCancel();
      }}
    >
      <AlertDialog.Portal>
        <AlertDialog.Overlay
          className='fixed inset-0 z-50 bg-black/30'
          onClick={(e) => {
            e.stopPropagation();
          }}
        />
        <AlertDialog.Content
          className='fixed top-1/2 left-1/2 z-50 w-[90vw] max-w-md -translate-x-1/2 -translate-y-1/2 rounded-lg border bg-white p-6 shadow-lg'
          onClick={(e) => {
            e.stopPropagation();
          }}
        >
          <AlertDialog.Title className='mb-4 text-xl font-bold'>
            {title}
          </AlertDialog.Title>
          <AlertDialog.Description className='mb-6 text-gray-600'>
            {message}
          </AlertDialog.Description>
          <div className='flex justify-end gap-3'>
            <AlertDialog.Cancel asChild>
              <button
                className='rounded-md border px-4 py-2 hover:bg-gray-50'
                onClick={(e) => {
                  e.stopPropagation();
                }}
              >
                {cancelText}
              </button>
            </AlertDialog.Cancel>
            <AlertDialog.Action asChild>
              <button
                className='rounded-md bg-red-600 px-4 py-2 text-white hover:bg-red-700'
                onClick={(e) => {
                  e.stopPropagation();
                  onConfirm();
                }}
              >
                {confirmText}
              </button>
            </AlertDialog.Action>
          </div>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  );
}
