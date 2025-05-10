import SignupForm from '@/app/ui/signup-form';
import { Suspense } from 'react';

export default async function SignupPage(props: {
  searchParams?: Promise<{
    accountType: string;
  }>;
}) {
  const searchParams = await props.searchParams;
  const baseAccountType = searchParams?.accountType || 'P';

  return (
    <main className="flex items-center justify-center min-h-screen">
      <div className="relative mx-auto flex w-full max-w-[400px] flex-col space-y-2.5 p-4">
        <Suspense>
          <SignupForm baseAccountType={baseAccountType} />
        </Suspense>
      </div>
    </main>
  );
}
