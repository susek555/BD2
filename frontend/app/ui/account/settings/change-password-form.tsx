import { KeyIcon } from '@heroicons/react/24/outline';
import { Button } from '../../(common)/button';
import InputField from '../../(common)/input-field';

export default function ChangePasswordForm() {
  return (
    <form className='space-y-3'>
      <InputField
        id='currentPassword'
        name='currentPassword'
        label='Current Password'
        placeholder='Enter your current password'
        type='password'
        defaultValue={''}
        icon={KeyIcon}
        errors={[]}
        required
      />
      <InputField
        id='newPassword'
        name='newPassword'
        label='New Password'
        placeholder='Enter your new password'
        type='password'
        defaultValue={''}
        icon={KeyIcon}
        errors={[]}
        required
      />
      <InputField
        id='confirmNewPassword'
        name='confirmNewPassword'
        label='Confirm New Password'
        placeholder='Confirm your new password'
        type='password'
        defaultValue={''}
        icon={KeyIcon}
        errors={[]}
        required
      />
      <div className='flex justify-end'>
        <Button type='submit' className='mt-2 px-4'>
          Change Password
        </Button>
      </div>
    </form>
  );
}
