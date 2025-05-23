import { cachedSessionData } from '@/app/lib/data/account/data';
import ChangePasswordForm from '@/app/ui/account/settings/change-password-form';
import ChangePersonalDataForm from '@/app/ui/account/settings/change-personal-data-form';

export default async function SettingsPage() {
  const userProfile = await cachedSessionData();
  return (
    <div className='p-6'>
      <div className='grid grid-cols-1 gap-6 md:grid-cols-2'>
        <div className='pl-20 pr-2'>
          <h2 className='mb-4 text-xl font-semibold'>Edit Profile</h2>
          <ChangePersonalDataForm userProfile={userProfile} />
        </div>
        <div className='pl-2 pr-20'>
          <h2 className='mb-4 text-xl font-semibold'>Change Password</h2>
          <ChangePasswordForm />
        </div>
      </div>
    </div>
  );
}
