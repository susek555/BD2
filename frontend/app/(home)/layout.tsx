import TopBar from "@/app/ui/(home)/topbar";
import { captureRejectionSymbol } from "events";

export default function Layout({children}: {children: React.ReactNode}) {
    return (
      <div className="flex flex-grow flex-col">
        <div className="h-full flex-none md:h-10">
          <TopBar />
        </div>
        {children}
      </div>
    )
}