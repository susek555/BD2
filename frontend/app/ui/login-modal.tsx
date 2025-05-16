'use client';

import { XMarkIcon } from '@heroicons/react/20/solid';
import * as Dialog from '@radix-ui/react-dialog';
import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import { useState } from 'react';
import LoginForm from './login-form';

interface LoginDialogProps {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  onLoginSuccess?: () => void;
}

export default function LoginModal({
  open,
  onOpenChange,
  onLoginSuccess,
}: LoginDialogProps) {
  const [dialogOpen, setDialogOpen] = useState(open || false);

  const isControlled = open !== undefined && onOpenChange !== undefined;
  const isOpen = isControlled ? open : dialogOpen;
  const setIsOpen = isControlled ? onOpenChange : setDialogOpen;

  const handleLoginSuccess = () => {
    setIsOpen(false);
    if (onLoginSuccess) {
      onLoginSuccess();
    }
  };

  return (
    <Dialog.Root open={isOpen} onOpenChange={setIsOpen}>
      <Dialog.Portal>
        <VisuallyHidden>
          <Dialog.Title>Login required</Dialog.Title>
          <Dialog.Description>
            You need to be logged in to perform this action
          </Dialog.Description>
        </VisuallyHidden>
        <Dialog.Overlay className='animate-fade-in fixed inset-0 z-40 bg-black/50' />
        <Dialog.Content
          className='animate-slide-up fixed top-1/2 left-1/2 z-50 max-h-[90vh] w-[90vw] max-w-md -translate-x-1/2 -translate-y-1/2 overflow-auto rounded-lg bg-white shadow-lg'
          onClick={(e) => {
            e.stopPropagation();
          }}
        >
          <div className='relative'>
            <div>
              <LoginForm onLoginSuccess={handleLoginSuccess} />
            </div>

            <Dialog.Close asChild>
              <button
                className='absolute top-4 right-4 rounded-full p-1 hover:bg-gray-100'
                aria-label='Close'
              >
                <XMarkIcon className='h-6 w-6' />
              </button>
            </Dialog.Close>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
