import { UserCircleIcon } from "@heroicons/react/20/solid";

export default function UserDetails({ sellerName }: { sellerName: string }) {
    return (
        <div className="flex flex-col gap-4 w-full md:w-120 h-full md:h-40 border-gray-300 border">
            <div className="flex justify-center items-center flex-col h-full gap-2">
                <p className="text-2xl">Seller</p>
                <a
                    href={`/profile/${sellerName}`}
                    // TODO - add link to user profile
                    className="flex flex-row justify-center md:w-100 items-center bg-gray-100 p-4 rounded-md hover:bg-gray-200 transition"
                >
                    <div className="flex flex-row gap-2 items-center">
                        <div className="w-8 h-8 rounded-full bg-blue-500 flex justify-center items-center">
                            <UserCircleIcon className="w-8 h-8 text-white" />
                        </div>
                        <p className="font-bold text-2xl">{sellerName}</p>
                    </div>
                </a>
            </div>
        </div>
    );
}