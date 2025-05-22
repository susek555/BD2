import { cachedSessionData } from '@/app/lib/data/account/data';
import ChangePersonalDataForm from '@/app/ui/account/settings/change-personal-data-form';

export default async function SettingsPage() {
  const userProfile = await cachedSessionData();
  return (
    <div className='p-6'>
      <h2 className='mb-4 text-xl font-semibold text-gray-900'>
        Account Settings
      </h2>
      <ChangePersonalDataForm userProfile={userProfile} />
    </div>
  );
}
