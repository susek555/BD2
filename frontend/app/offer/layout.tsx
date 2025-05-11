import { OfferTopBar } from "@/app/ui/(topbar)/topbar";

export default function Layout({children}: {children: React.ReactNode}) {
    return (
      <div className="flex flex-grow flex-col">
        <div className="h-full flex-none md:h-10">
          <OfferTopBar />
        </div>
        {children}
      </div>
    )
}