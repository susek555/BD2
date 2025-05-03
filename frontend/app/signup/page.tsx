import SignupForm from '@/app/ui/signup-form';
import { Suspense } from 'react';

export default function SignupPage() {
  return (
    <main className="flex items-center justify-center min-h-screen">
      <div className="relative mx-auto flex w-full max-w-[400px] flex-col space-y-2.5 p-4">
        <Suspense>
          <SignupForm />
        </Suspense>
      </div>
    </main>
  );
}
