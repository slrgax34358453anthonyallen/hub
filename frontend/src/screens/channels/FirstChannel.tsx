import { Link } from "react-router-dom";
import Loading from "src/components/Loading";
import { AlbyIcon } from "src/components/icons/Alby";
import { MIN_ALBY_BALANCE } from "src/constants";
import { useAlbyBalance } from "src/hooks/useAlbyBalance";
import { useAlbyMe } from "src/hooks/useAlbyMe";
import { useInfo } from "src/hooks/useInfo";

export default function FirstChannel() {
  const { data: info } = useInfo();
  const { data: albyMe } = useAlbyMe();
  const { data: albyBalance } = useAlbyBalance();

  if (!info) {
    return <Loading />;
  }

  return (
    <div className="flex flex-col justify-center items-center gap-4">
      <p className="max-w-lg">
        Your Alby has grown up now. You have your own node and you also need
        your own channels to send and receive payments on the lightning network.
      </p>

      {!albyMe && (
        <>
          <p>
            If you have funds on your Alby account you can use them to open your
            first channel.
          </p>

          <Link to={info.albyAuthUrl}>
            <button
              className="font-body flex h-10 w-56 items-center justify-center gap-2 rounded-md font-bold text-black shadow transition-all hover:brightness-90 active:scale-95"
              style={{
                background:
                  "linear-gradient(180deg, #FFDE6E 63.72%, #F8C455 95.24%)",
              }}
              type="button"
            >
              <AlbyIcon className="w-6 h-6" />
              <span className="mr-2">Connect with Alby</span>
            </button>
          </Link>
        </>
      )}
      {albyMe && albyBalance && (
        <div className="mt-8 border-2 p-8 rounded-lg flex flex-col justify-center items-center border-yellow-300 bg-yellow-100">
          <p className="mb-8">
            Logged in as{" "}
            <span className="font-bold">{albyMe.lightning_address}</span>
          </p>

          {albyBalance.sats >= MIN_ALBY_BALANCE && (
            <>
              <Link to="/channels/migrate-alby">
                <button className="bg-yellow-400 border-8 rounded-lg border-yellow-500 p-4 shadow-lg font-mono text-lg font-black">
                  Migrate Funds 🚀
                </button>
              </Link>
              <p className="text-sm italic mt-4">
                You have {albyBalance.sats} sats to migrate
              </p>
            </>
          )}
          {albyBalance.sats < MIN_ALBY_BALANCE && (
            <>
              <p>
                You don't have enough sats in your Alby account to open a
                channel.
              </p>
            </>
          )}
        </div>
      )}

      <p className="mt-8">-- or --</p>

      <Link to="/channels/new" className="mt-8 text-purple-500">
        Fund & Open Channel
      </Link>
    </div>
  );
}
