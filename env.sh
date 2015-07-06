if [ "$INCLUDE_QBOX_BASE" = "" ]; then
	source $QBOXROOT/base/env.sh
fi

export GOPATH=$QBOXROOT/boltfs:$GOPATH
export PATH=$PATH:$QBOXROOT/boltfs/bin
export INCLUDE_QBOX_BOLTFS=1

