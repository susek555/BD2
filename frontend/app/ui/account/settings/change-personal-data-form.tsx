'use client';

import { revalidateSesion } from '@/app/actions/revalidate-session';
import { updateProfile } from '@/app/actions/update-profile';
import { PersonalDataFormState } from '@/app/lib/definitions/account-settings';
import { UserProfile } from '@/app/lib/definitions/user';
import {
  AtSymbolIcon,
  BuildingOfficeIcon,
  IdentificationIcon,
  UserCircleIcon,
  UserIcon,
} from '@heroicons/react/24/outline';
import { useSession } from 'next-auth/react';
import { useActionState } from 'react';
import { Button } from '../../(common)/button';
import InputField from '../../(common)/input-field';

interface PersonalDataFormProps {
  userProfile: UserProfile;
}

export default function ChangePersonalDataForm({
  userProfile,
}: PersonalDataFormProps) {
  const { update: updateSession } = useSession();

  const initialState: PersonalDataFormState = {
    errors: {},
    values: {},
  };

  const wrappedUpdateProfile = async (
    state: PersonalDataFormState,
    formData: FormData,
  ) => {
    const username = formData.get('username') as string;
    const email = formData.get('email') as string;

    let formUnchanged =
      username === userProfile.username && email === userProfile.email;

    if (userProfile.selector === 'P') {
      const personName = formData.get('person_name') as string;
      const personSurname = formData.get('person_surname') as string;
      formUnchanged =
        formUnchanged &&
        personName === (userProfile.personName || '') &&
        personSurname === (userProfile.personSurname || '');
    } else if (userProfile.selector === 'C') {
      const companyName = formData.get('company_name') as string;
      const companyNip = formData.get('company_nip') as string;
      formUnchanged =
        formUnchanged &&
        companyName === (userProfile.companyName || '') &&
        companyNip === (userProfile.companyNip || '');
    }

    if (formUnchanged) {
      return state;
    }

    const updateResult = await updateProfile(state, formData, userProfile.id);

    if (!updateResult.state.errors && updateResult.newUserData) {
      await updateSession(updateResult.newUserData);
      await revalidateSesion('/account', 'layout');
    }

    return updateResult.state;
  };

  const [state, action] = useActionState(wrappedUpdateProfile, initialState);

  const handleSubmit = (formData: FormData) => {
    return action(formData);
  };

  return (
    <form className='space-y-3' action={handleSubmit}>
      <input
        type='hidden'
        name='selector'
        defaultValue={userProfile.selector}
      />
      <InputField
        id='username'
        name='username'
        label='Username'
        placeholder='Enter your username'
        defaultValue={state?.values?.username || userProfile.username}
        icon={IdentificationIcon}
        errors={state?.errors?.username || []}
        required
      />
      <InputField
        id='email'
        name='email'
        label='Email'
        placeholder='Enter your email address'
        defaultValue={state?.values?.email || userProfile.email}
        icon={AtSymbolIcon}
        errors={state?.errors?.email || []}
        required
      />
      {userProfile.selector === 'P' && (
        <>
          <InputField
            id='person_name'
            name='person_name'
            label='Name'
            placeholder='Enter your name'
            defaultValue={
              state?.values?.personName || userProfile.personName || ''
            }
            icon={UserIcon}
            errors={state?.errors?.personName || []}
            required
          />
          <InputField
            id='person_surname'
            name='person_surname'
            label='Surname'
            placeholder='Enter your surname'
            defaultValue={
              state?.values?.personSurname || userProfile.personSurname || ''
            }
            icon={UserCircleIcon}
            errors={state?.errors?.personSurname || []}
            required
          />
        </>
      )}
      {userProfile.selector === 'C' && (
        <>
          <InputField
            id='company_name'
            name='company_name'
            label='Company Name'
            placeholder='Enter your company name'
            defaultValue={
              state?.values?.companyName || userProfile.companyName || ''
            }
            icon={BuildingOfficeIcon}
            errors={state?.errors?.companyName || []}
            required
          />
          <InputField
            id='company_nip'
            name='company_nip'
            label='NIP (Tax ID)'
            placeholder='Enter your NIP'
            defaultValue={
              state?.values?.companyNip || userProfile.companyNip || ''
            }
            icon={BuildingOfficeIcon}
            errors={state?.errors?.companyNip || []}
            required
          />
        </>
      )}
      <div className='flex justify-end'>
        <Button type='submit' className='mt-2 px-4'>
          Save
        </Button>
      </div>
    </form>
  );
}
