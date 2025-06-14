import {
  fetchAverageRating,
  fetchRatingDistribution,
} from '@/app/lib/data/reviews/data';
import { UserProfile } from '@/app/lib/definitions/user';
import { UserCircleIcon } from '@heroicons/react/24/outline';
import RatingDisplay from './rating-display';

interface ProfileInfoProps {
  user: UserProfile;
}

export default async function ProfileInfo({ user }: ProfileInfoProps) {
  const rating = await fetchAverageRating(user.userId);
  const distribution = await fetchRatingDistribution(user.userId);

  return (
    <div className='relative rounded-lg bg-white p-6 shadow'>
      <div className='absolute top-6 right-6'>
        <RatingDisplay distribution={distribution} rating={rating} />
      </div>

      <div className='flex flex-col items-center gap-4 md:flex-row md:items-start'>
        <div className='flex h-24 w-24 items-center justify-center rounded-full bg-gray-200'>
          <UserCircleIcon className='h-20 w-20 text-gray-400' />
        </div>
        <div className='flex flex-1 flex-col gap-2 text-center md:text-left'>
          <h1 className='text-2xl font-bold text-gray-900'>
            {user.personName && user.personSurname
              ? `${user.personName} ${user.personSurname}`
              : user.companyName || user.username}
          </h1>
          <p className='text-gray-600'>{user.username}</p>
          <p className='text-gray-600'>{user.email}</p>
          {user.companyNip && (
            <p className='text-gray-600'>NIP: {user.companyNip}</p>
          )}
        </div>
      </div>
    </div>
  );
}
