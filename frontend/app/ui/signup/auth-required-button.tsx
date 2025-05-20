'use client';

import React, { useState } from 'react';
import LoginModal from './login-modal';

interface AuthRequiredButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  isLoggedIn: boolean;
  onClick: () => void;
  children: React.ReactNode;
  dialogTitle?: string;
  dialogDescription?: string;
  dialogContent?: React.ReactNode;
}

export function AuthRequiredButton({
  isLoggedIn,
  onClick,
  children,
  ...buttonProps
}: AuthRequiredButtonProps) {
  const [dialogOpen, setDialogOpen] = useState(false);

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    e.preventDefault();

    if (isLoggedIn) {
      onClick();
    } else {
      setDialogOpen(true);
    }
  };

  return (
    <>
      <button {...buttonProps} onClick={handleClick}>
        {children}
      </button>

      {!isLoggedIn && (
        <LoginModal
          open={dialogOpen}
          onOpenChange={setDialogOpen}
          onLoginSuccess={onClick}
        />
      )}
    </>
  );
}
