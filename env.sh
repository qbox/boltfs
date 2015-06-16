if [ "$INCLUDE_QBOX_BASE" = "" ]; then
	source $QBOXROOT/base/env.sh
fi

export GOPATH=$QBOXROOT/qbs:$GOPATH
export PATH=$PATH:$QBOXROOT/qbs/bin
export INCLUDE_QBOX_QBS=1

