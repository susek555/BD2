'use client';

import { changePasswordAction } from '@/app/actions/change-password';
import { ChangePasswordFormState } from '@/app/lib/definitions/account-settings';
import { UserProfile } from '@/app/lib/definitions/user';
import { KeyIcon } from '@heroicons/react/24/outline';
import { useActionState } from 'react';
import { Button } from '../../(common)/button';
import InputField from '../../(common)/input-field';

interface ChangePasswordFormProps {
  userProfile: UserProfile;
}

export default function ChangePasswordForm({
  userProfile,
}: ChangePasswordFormProps) {
  const initialState: ChangePasswordFormState = {
    errors: {},
  };

  const wrappedChangePassword = async (
    state: ChangePasswordFormState,
    formData: FormData,
  ) => {
    return changePasswordAction(state, formData, userProfile.userId);
  };

  const [state, action] = useActionState(wrappedChangePassword, initialState);

  return (
    <form className='space-y-3' action={action}>
      <InputField
        id='current_password'
        name='current_password'
        label='Current Password'
        placeholder='Enter your current password'
        type='password'
        defaultValue={''}
        icon={KeyIcon}
        errors={state?.errors?.current_password || []}
        required
      />
      <InputField
        id='new_password'
        name='new_password'
        label='New Password'
        placeholder='Enter your new password'
        type='password'
        defaultValue={''}
        icon={KeyIcon}
        errors={state?.errors?.new_password || []}
        required
      />
      <InputField
        id='confirm_new_password'
        name='confirm_new_password'
        label='Confirm New Password'
        placeholder='Confirm your new password'
        type='password'
        defaultValue={''}
        icon={KeyIcon}
        errors={state?.errors?.confirm_new_password || []}
        required
      />
      <div className='flex justify-end'>
        <Button type='submit' className='mt-2 px-4'>
          Change Password
        </Button>
      </div>
      {state.success && (
        <p className='text-right text-sm text-green-600'>
          Password changed successfully
        </p>
      )}
    </form>
  );
}
