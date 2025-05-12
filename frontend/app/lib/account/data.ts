import { getServerSession } from 'next-auth';
import { authConfig } from '../authConfig';
import { UserProfile } from '../definitions';

export async function getSessionData(): Promise<UserProfile> {
  const user = (await getServerSession(authConfig))?.user;

  if (!user) {
    throw new Error('No session found');
  }

  const userProfile: UserProfile = {
    username: user.username,
    email: user.email,
  };

  if (user.selector === 'P') {
    userProfile.personName = user.personName;
    userProfile.personSurname = user.personSurname;
  } else if (user.selector === 'C') {
    userProfile.companyName = user.companyName;
    userProfile.companyNip = user.companyNip;
  }

  return userProfile;
}
