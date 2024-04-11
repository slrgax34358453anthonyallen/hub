import { ChevronRight, WalletMinimal } from "lucide-react";
import React from "react";
import { Link } from "react-router-dom";

import Container from "src/components/Container";
import TwoColumnLayoutHeader from "src/components/TwoColumnLayoutHeader";
import { Button } from "src/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle,
} from "src/components/ui/card";
import { useInfo } from "src/hooks/useInfo";

export function SetupWallet() {
  const { data: info } = useInfo();
  const [showOtherOptions, setShowOtherOptions] = React.useState(false);
  return (
    <>
      <Container>
        <div className="grid gap-5">
          <TwoColumnLayoutHeader
            title="Create Your Wallet"
            description="Alby Hub requires a wallet to connect to your apps."
          />
          {info?.backendType && (
            <>
              <Link to="/setup/finish">
                <Card>
                  <CardContent className="pt-6">
                    <div className="flex gap-3 items-center">
                      <WalletMinimal className="w-8 h-8 p-1 text-green-500" />
                      <div className="flex-grow">
                        <CardTitle>{info.backendType} Wallet</CardTitle>
                        <CardDescription>
                          Connect to preconfigured {info.backendType} Wallet
                        </CardDescription>
                      </div>
                      <div className="flex-shrink-0 flex justify-end ">
                        <ChevronRight className="w-8 h-8 text-muted-foreground" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </Link>

              {!showOtherOptions && (
                <Button
                  className="mt-4"
                  onClick={() => setShowOtherOptions(true)}
                >
                  See other options
                </Button>
              )}
            </>
          )}

          {(showOtherOptions || !info?.backendType) && (
            <div className="grid gap-1.5 justify-center text-center">
              <Link to="/setup/node?wallet=new">
                <Button>New Wallet</Button>
              </Link>
              <Link to="/setup/node?wallet=import">
                <Button variant="link">Import Existing Wallet</Button>
              </Link>
            </div>
          )}
        </div>
      </Container>
    </>
  );
}
